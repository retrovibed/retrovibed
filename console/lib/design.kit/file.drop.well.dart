import 'dart:io';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:desktop_drop/desktop_drop.dart';
import 'package:file_selector/file_selector.dart';

class FilesEvent {
  final List<DropItemFile> files;
  const FilesEvent({required this.files});
}

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
  final EdgeInsets? margin;
  final Future<Widget?> Function(
    FilesEvent i, {
    ValueNotifier<int>? progress,
  })
  onDropped;
  final List<String> mimetypes;
  final List<String> extensions;
  final Widget help;
  final String? tooltip;

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
    this.margin,
    this.help = ds.HelpScope.None,
    this.tooltip,
  });

  static Future<FilesEvent> files({
    List<String> mimetypes = const [],
    List<String> extensions = const [],
  }) {
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
      ValueNotifier<int>? progress,
    })
    onDropped, {
    Key? key,
    List<String> mimetypes = const [],
    List<String> extensions = const [],
    IconData icon = Icons.file_upload_outlined,
    Function()? onTap,
    Widget help = ds.HelpScope.None,
    String? tooltip,
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
    );
  }

  @override
  _FileDropWell createState() => _FileDropWell();
}

class _FileDropWell extends State<FileDropWell> {
  final ValueNotifier<int> _progress = ValueNotifier(0);
  int _total = 0;
  bool _dragging = false;
  bool _loading = false;

  @override
  void initState() {
    super.initState();
    _progress.addListener(() {
      setState(() {});
    });
  }

  @override
  void dispose() {
    _progress.dispose();
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
            final total = resolved.files.fold<int>(0, (acc, f) => acc + File(f.path).lengthSync());
            setState(() {
              _progress.value = 0;
              _total = total;
            });
            return widget.onDropped(resolved, progress: _progress);
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
            Future.wait(
                  evt.files.map((c) {
                    return c.openRead(0, mimex.defaultMagicNumbersMaxLength).first.then((v) => v.toList()).then((bits) {
                      return new DropItemFile(
                        c.path,
                        name: c.name,
                        mimeType: mimex.fromFile(c.name, magicbits: bits).toString(),
                      );
                    });
                  }),
                )
                .then((files) {
                  final resolved = FilesEvent(files: files);
                  widget.onDropped(resolved).whenComplete(() {
                    setState(() {
                      _loading = false;
                    });
                  });
                })
                .catchError((cause) {
                  print("failed to open file dialog ${cause}");
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
            child: ds.LoadingIconButton(
              onPressed: onPress,
              icon: widget.child,
              disabled: _loading,
              tooltip: widget.tooltip,
              value: _progress.value / _total,
            ),
          ),
        ),
      ),
      widget.help,
    );
  }
}
