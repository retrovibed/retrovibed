import 'dart:async';
import 'dart:io';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:desktop_drop/desktop_drop.dart';
import 'package:file_selector/file_selector.dart';

class FilesEvent {
  // DropItemFile for files, DropItemDirectory for directories.
  final List<DropItem> files;
  const FilesEvent({required this.files});
}

// a well that only accepts directories, e.g. mimetypes: mimex.folders.
bool _foldersonly(List<String> mimetypes) => mimetypes.isNotEmpty && mimetypes.every((m) => m == mimex.directory);

class FileDropWell extends StatefulWidget {
  static Widget textual(String text) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.max,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [Icon(Icons.filter_rounded), SelectableText(text)],
      ),
    );
  }

  final Widget child;
  final Widget? loading;
  final Function()? onTap;
  final EdgeInsets margin;
  final EdgeInsets padding;
  final Future<Widget?> Function(
    FilesEvent i, {
    StreamSink<httpx.UploadProgress>? progress,
  })
  onDropped;
  final List<String> mimetypes;
  final List<String> extensions;
  final Widget help;
  final String? tooltip;
  final OutlinedBorder? shape;

  const FileDropWell(
    this.onDropped, {
    super.key,
    this.child = const Center(
      child: Column(
        mainAxisSize: MainAxisSize.max,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.filter_rounded),
          SelectableText("Drop Files to add them to your library."),
        ],
      ),
    ),
    this.mimetypes = const [],
    this.extensions = const [],
    this.onTap,
    this.loading,
    this.margin = EdgeInsets.zero,
    this.padding = EdgeInsets.zero,
    this.help = ds.HelpScope.None,
    this.tooltip,
    this.shape,
  });

  // resolves dropped items: directories are detected directly rather than sniffed (desktop_drop
  // only reports directories on macOS) and are kept only when the mimetypes accept them. a
  // folders only well discards files, every other well passes files through unfiltered.
  static Future<FilesEvent> resolve(List<XFile> items, {List<String> mimetypes = const []}) {
    final directories = mimetypes.contains(mimex.directory);
    final foldersonly = _foldersonly(mimetypes);

    final resolved = items.map((c) {
      if (FileSystemEntity.isDirectorySync(c.path)) {
        if (!directories) return Future<DropItem?>.value(null);
        return Future<DropItem?>.value(DropItemDirectory(c.path, const [], name: c.name, mimeType: mimex.directory));
      }

      if (foldersonly) return Future<DropItem?>.value(null);

      return c.openRead(0, mimex.defaultMagicNumbersMaxLength).first.then((v) => v.toList()).then((bits) {
        return DropItemFile(
          c.path,
          name: c.name,
          mimeType: mimex.fromFile(c.name, magicbits: bits).toString(),
        );
      });
    });

    return Future.wait(resolved).then((files) => FilesEvent(files: files.nonNulls.toList()));
  }

  static Future<FilesEvent> files({
    List<String> mimetypes = const [],
    List<String> extensions = const [],
  }) {
    if (_foldersonly(mimetypes)) {
      return getDirectoryPath().then((path) {
        if (path == null) return const FilesEvent(files: []);
        return FilesEvent(
          files: [
            DropItemDirectory(path, const [], name: path.split(Platform.pathSeparator).last, mimeType: mimex.directory),
          ],
        );
      });
    }

    final XTypeGroup filter = XTypeGroup(
      label: "Select File(s)",
      extensions: extensions,
      mimeTypes: mimetypes,
    );

    return openFiles(acceptedTypeGroups: [filter]).then((files) {
      final eventfiles = files.map((f) {
        final fh = File(f.path);
        return fh.openSync().read(mimex.defaultMagicNumbersMaxLength).then((v) => v.toList()).then((bits) {
          return DropItemFile(
            f.path,
            name: f.name,
            mimeType: mimex.fromFile(f.name, magicbits: bits).toString(),
          );
        });
      }).toList();

      return Future.wait(eventfiles).then((files) => FilesEvent(files: files));
    });
  }

  factory FileDropWell.icon(
    Future<Widget?> Function(
      FilesEvent i, {
      StreamSink<httpx.UploadProgress>? progress,
    })
    onDropped, {
    Key? key,
    List<String> mimetypes = const [],
    List<String> extensions = const [],
    IconData icon = Icons.file_upload_outlined,
    Function()? onTap,
    Widget help = ds.HelpScope.None,
    String? tooltip,
    OutlinedBorder? shape,
  }) {
    return FileDropWell(
      onDropped,
      key: key,
      onTap: onTap,
      child: Icon(icon, size: 24.0),
      loading: ds.Loading.Sized(width: 24.0, height: 24.0),
      mimetypes: mimetypes,
      extensions: extensions,
      help: help,
      tooltip: tooltip,
      shape: shape,
    );
  }

  @override
  _FileDropWell createState() => _FileDropWell();
}

class _FileDropWell extends State<FileDropWell> {
  final StreamController<httpx.UploadProgress> _progress = StreamController<httpx.UploadProgress>();
  final Map<String, int> _uploaded = {};
  int _total = 0;
  bool _dragging = false;
  bool _loading = false;

  @override
  void initState() {
    super.initState();
    _progress.stream.listen((event) {
      setState(() {
        _uploaded[event.$1] = event.$4;
      });
    });
  }

  @override
  void dispose() {
    _progress.close();
    super.dispose();
  }

  @override
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    Future<void> onPress() {
      return FileDropWell.files(mimetypes: widget.mimetypes, extensions: widget.extensions)
          .then((resolved) {
            if (resolved.files.isEmpty) return null;
            // directories have no length of their own.
            final total = resolved.files.whereType<DropItemFile>().fold<int>(
              0,
              (acc, f) => acc + File(f.path).lengthSync(),
            );
            setState(() {
              _uploaded.clear();
              _total = total;
            });
            return widget.onDropped(resolved, progress: _progress.sink);
          })
          .catchError((cause) {
            return ds.Error.unknown(cause);
          });
    }

    return ds.Help(
      Material(
        // Ensure Material doesn't block underlying colors
        color: Colors.transparent,
        child: DropTarget(
          onDragDone: (evt) {
            setState(() {
              _loading = true;
            });
            FileDropWell.resolve(evt.files, mimetypes: widget.mimetypes)
                .then((resolved) {
                  if (resolved.files.isEmpty) return null;
                  return widget.onDropped(resolved);
                })
                .catchError((cause) {
                  print("failed to open file dialog ${cause}");
                  return null;
                })
                .whenComplete(() {
                  setState(() {
                    _loading = false;
                  });
                });
          },
          onDragEntered: (detail) {
            setState(() {
              _dragging = true;
            });
          },
          onDragExited: (detail) {
            setState(() {
              _dragging = false;
            });
          },
          child: Container(
            color: _dragging ? theme.highlightColor : null,
            margin: widget.margin,
            padding: widget.padding,
            child: ds.LoadingIconButton(
              onPressed: onPress,
              icon: widget.child,
              disabled: _loading,
              tooltip: widget.tooltip,
              value: _uploaded.values.fold<int>(0, (a, b) => a + b) / _total,
              shape: widget.shape,
            ),
          ),
        ),
      ),
      widget.help,
    );
  }
}
