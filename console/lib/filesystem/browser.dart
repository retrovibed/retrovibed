import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/file.drop.well.dart';
import 'package:retrovibed/filesystem/api.dart' as api;
import 'package:retrovibed/library/dropdown.upload.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'directory.create.dart';
import 'row.dart';

// browses the library as a tree. this is a sibling of the library view rather than a
// variation on it: the two share the Media row and nothing else, because the library grid
// is flat by definition and never shows a directory.
class FilesystemBrowser extends StatefulWidget {
  final api.FnFilesystemSearch search;
  final api.FnFilesystemCreate create;
  final api.FnFilesystemDelete remove;
  final media.FnUploadRequest upload;
  final TextEditingController? controller;
  final FocusNode? focus;
  final ValueNotifier<media.SearchMode> mode;
  final void Function(media.SearchMode) onModeChanged;

  const FilesystemBrowser({
    super.key,
    required this.mode,
    required this.onModeChanged,
    this.search = api.filesystem.search,
    this.create = api.filesystem.create,
    this.remove = api.filesystem.delete,
    this.upload = media.media.upload,
    this.controller,
    this.focus,
  });

  @override
  State<StatefulWidget> createState() => _FilesystemBrowser();
}

class _FilesystemBrowser extends State<FilesystemBrowser> {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  String _highlighted = "";
  api.FilesystemSearchResponse _res = api.filesystem.response(
    next: api.filesystem.request(limit: 32),
  );

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  String get directory => _res.next.directoryId;

  // the directory holding the one being listed. the daemon returns the path root first, so
  // the entry before the last is the parent; a single entry means the parent is the root.
  String get ancestor {
    final path = _res.breadcrumb;
    return path.length < 2 ? uuidx.min() : path[path.length - 2].id;
  }

  String get location {
    final path = _res.breadcrumb.map((v) => v.description).join("/");
    return path.isEmpty ? "search the library" : "search in ${path}";
  }

  Future<void> refresh(api.FilesystemSearchRequest req) {
    return widget
        .search(req, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _res = v;
            _loading = false;
          });

          widget.focus?.requestFocus();
          ds.textediting.refocus(widget.controller);
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unauthorized(cause, onTap: reseterr);
            _loading = false;
          });
        }, test: httpx.ErrorsTest.unauthorized)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: reseterr);
            _loading = false;
          });
        });
  }

  void navigate(String id) {
    setState(() {
      _loading = true;
      _highlighted = "";
      _res.next
        ..directoryId = id
        ..offset = ds.Int64(0);
    });
    refresh(_res.next);
  }

  @override
  void initState() {
    super.initState();
    _res.next.query = widget.controller?.text ?? "";
    ds.postframe(() => refresh(_res.next));
  }

  // moving up is an entry at the head of the listing rather than a breadcrumb bar, so it
  // costs no chrome and renders through the same row widget as everything else.
  List<media.Media> get items {
    if (uuidx.fromString(directory) == uuidx.fromString(uuidx.min())) return _res.items;
    return [media.Media(id: ancestor, description: "..", mimetype: mimex.directory), ..._res.items];
  }

  Widget row(media.Media v) {
    if (v.mimetype == mimex.directory) {
      return FilesystemRow(
        current: v,
        highlighted: v.id == _highlighted,
        trailing: [
          Visibility(
            visible: v.description != "..",
            child: IconButton(
              icon: const Icon(Icons.delete_outline),
              onPressed: () => confirmremove(v),
            ),
          ),
        ],
        onTap: () => Future.sync(() => navigate(v.id)),
      );
    }

    return FilesystemRow(
      current: v,
      highlighted: v.id == _highlighted,
      trailing: [media.ButtonShare(current: v)],
      // playback owns audio and video; everything else the reader can inspect without
      // pulling the whole file down first.
      onTap: media.PlayAction(context, v, media.media.response()) ?? preview(v),
    );
  }

  Future<void> Function() preview(media.Media v) {
    final constraints = ds.Defaults.of(context).modal(context);

    return () => Future.sync(
      () => ds.modals.push(
        context,
        ds.Card(
          constraints: constraints,
          leading: [ds.Heading(Text(v.description))],
          trailing: [
            IconButton(icon: const Icon(Icons.close), onPressed: () => ds.modals.push(context, null)),
          ],
          media.Preview(current: v),
        ),
      ),
    );
  }

  // deleting a directory deletes what it holds, which is not recoverable from this screen,
  // so the user is told before it happens rather than after.
  void confirmremove(media.Media v) {
    final modal = ds.modals.of(context);
    modal?.push(
      ds.Confirmation.yesNo(
        content: Text(
          "Delete ${v.description}? Everything inside it is removed from your library too.",
        ),
        onCancel: (_) => modal.push(null),
        onConfirm: (_) {
          modal.push(null);
          setState(() {
            _loading = true;
          });

          httpx
              .withRetry(
                () => widget.remove(v.id, options: [authn.request(authn.AuthzCache.meta(context))]),
              )
              .then((_) => refresh(_res.next))
              .catchError((cause) {
                setState(() {
                  _cause = ds.Error.unknown(cause, onTap: reseterr);
                  _loading = false;
                });
              });
        },
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final upload = (FilesEvent v, {ValueNotifier<int>? progress}) {
      setState(() {
        _loading = true;
      });

      final multiparts = v.files.map((c) => media.media.uploadable(c.path, c.name, c.mimeType!));

      return Future.microtask(() {
        return Future.wait(
              multiparts.map(
                (fv) => fv.then(
                  (v) => widget.upload((req) {
                    // files dropped onto the listing belong to the directory on screen.
                    req.fields["directory_id"] = directory;
                    req.files.add(v);
                    return req;
                  }),
                ),
              ),
            )
            .then((_) => refresh(_res.next))
            .then((_) => ds.NullWidget)
            .catchError((cause) => ds.Error.unknown(cause, onTap: reseterr));
      }).whenComplete(() => setState(() => _loading = false));
    };

    return ds.Table(
      loading: _loading,
      cause: _cause,
      leading: ds.SearchTray(
        autofocus: defaults.desktop,
        decoration: InputDecoration(hintText: location),
        controller: widget.controller,
        focus: widget.focus,
        onSubmitted: (v) {
          setState(() {
            _res.next
              ..query = v
              ..offset = ds.Int64(0);
          });
          return refresh(_res.next);
        },
        next: (i) {
          setState(() {
            _res.next.offset = i;
          });
          refresh(_res.next);
        },
        current: _res.next.offset,
        empty: ds.Int64(_res.items.length) < _res.next.limit,
        leading: [
          ds.CompactingMenu.pinned(
            DropdownUpload(
              icon: Icon(mimex.icofolder),
              help: ds.Hint(const Text("create a directory, or switch to library or discover mode")),
              items: [
                PopupMenuItem<String>(
                  onTap: mkdir,
                  child: ListTile(leading: Icon(mimex.icofolder), title: const Text("New Directory")),
                ),
                const PopupMenuDivider(),
                media.SearchModeToggle(
                  mode: media.SearchMode.library,
                  current: widget.mode,
                  icon: Icons.video_library,
                  label: "Library",
                  onSelect: widget.onModeChanged,
                ),
                media.SearchModeToggle(
                  mode: media.SearchMode.discovery,
                  current: widget.mode,
                  icon: Icons.travel_explore,
                  label: "Discover",
                  onSelect: widget.onModeChanged,
                ),
                media.SearchModeToggle(
                  mode: media.SearchMode.downloads,
                  current: widget.mode,
                  icon: Icons.download,
                  label: "Downloads",
                  onSelect: widget.onModeChanged,
                ),
              ],
            ),
          ),
          ds.FileDropWell.icon(
            upload,
            help: ds.Hint(const Text("drag and drop files to add them to this directory")),
          ),
        ],
      ),
      children: items,
      ds.Table.expanded<media.Media>(row),
      empty: ds.FileDropWell(upload),
    );
  }

  void mkdir() {
    ds.modals.push(
      context,
      DirectoryCreate(
        parent: directory,
        create: widget.create,
        onCancel: () => ds.modals.push(context, null),
        onCreated: (created) {
          ds.modals.push(context, null);
          setState(() {
            _highlighted = created.id;
          });
          refresh(_res.next);
        },
      ),
    );
  }
}
