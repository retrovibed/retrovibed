import 'dart:async';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/file.drop.well.dart';
import 'package:retrovibed/filesystem/api.dart' as api;
import 'package:retrovibed/library/dropdown.nav.menu.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'directory.create.dart';
import 'details.dart';
import 'autoimport.search.dart';

// browses the library as a tree. this is a sibling of the library view rather than a
// variation on it: the two share the Media row and nothing else, because the library grid
// is flat by definition and never shows a directory.
class FilesystemBrowser extends StatefulWidget {
  final ValueNotifier<media.MediaSearchState> search;
  final api.FnFilesystemSearch apisearch;
  final api.FnFilesystemCreate apicreate;
  final api.FnFilesystemDelete apiremove;
  final media.FnUploadRequest apiupload;
  final TextEditingController? controller;
  final FocusNode? focus;
  final ValueNotifier<media.SearchMode> mode;
  final void Function(media.SearchMode) onModeChanged;

  const FilesystemBrowser({
    super.key,
    required this.search,
    required this.mode,
    required this.onModeChanged,
    this.apisearch = api.filesystem.search,
    this.apicreate = api.filesystem.create,
    this.apiremove = api.filesystem.delete,
    this.apiupload = media.media.upload,
    this.controller,
    this.focus,
  });

  @override
  State<StatefulWidget> createState() => _FilesystemBrowser();
}

class _FilesystemBrowser extends State<FilesystemBrowser> with ds.LoadingState {
  // the entry whose details are open; only the info button sets it.
  String _focused = "";
  // a view covering the browser completely (e.g. monitored directories), the browser stays
  // mounted beneath it so closing returns to the same directory.
  Widget _overlay = ds.Empty;
  api.FilesystemSearchResponse _res = api.filesystem.response(
    next: api.filesystem.request(limit: 32),
  );

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
        .apisearch(req, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _res = v;
            loading = false;
          });

          widget.focus?.requestFocus();
          ds.textediting.refocus(widget.controller);
        })
        .catchError((e) {
          setState(() {
            cause = ds.Error.unauthorized(e, onTap: reseterr);
            loading = false;
          });
        }, test: httpx.ErrorsTest.unauthorized)
        .catchError((e) {
          setState(() {
            cause = ds.Error.unknown(e, onTap: reseterr);
            loading = false;
          });
        });
  }

  // toggles the details for the entry.
  void Function() focus(media.Media v) {
    return () => setState(() => _focused = _focused == v.id ? "" : v.id);
  }

  void navigate(String id) {
    setState(() {
      loading = true;
      _focused = "";
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

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final upload = (FilesEvent v, {StreamSink<httpx.UploadProgress>? progress}) {
      final tracked = meta.UploadNode.of(context).progress;

      return Future.microtask(() {
        return Future.wait(
              v.files.map((c) {
                final abort = Completer<void>();
                return media.media
                    .uploadable(c.path, c.name, c.mimeType!, progress: tracked, abort: abort)
                    .then(
                      (v) => widget.apiupload(
                        (method, url) => http.AbortableMultipartRequest(method, url, abortTrigger: abort.future)
                          // files dropped onto the listing belong to the directory on screen.
                          ..fields["directory_id"] = directory
                          ..files.add(v),
                      ),
                    )
                    .catchError((_) => media.MediaUploadResponse(), test: httpx.ErrorsTest.aborted);
              }),
            )
            .then((_) => refresh(_res.next))
            .then((_) => ds.NullWidget)
            .catchError((e) => ds.Error.unknown(e, onTap: reseterr));
      });
    };

    final table = ds.Table(
      loading: loading,
      cause: cause,
      empty: ds.FileDropWell(
        upload,
        shape: RoundedRectangleBorder(borderRadius: defaults.borderRadius),
      ),
      leading: Column(
        mainAxisSize: MainAxisSize.min,
        verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
        children: [
          ds.SearchTray(
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
                DropdownNavMenu.options(
                  icon: Icon(mimex.icofolder),
                  help: ds.Hint(
                    const Text(
                      "filter by mimetype, create a directory, or switch between library, files, discover, and downloads mode",
                    ),
                  ),
                  search: widget.search,
                  mode: widget.mode,
                  onModeChanged: widget.onModeChanged,
                  options: [
                    PopupMenuItem<String>(
                      onTap: mkdir,
                      child: ListTile(leading: Icon(mimex.icofolder), title: const Text("New Folder")),
                    ),
                    PopupMenuItem<String>(
                      onTap: () => overlay(AutoimportSearch(onClose: () => overlay(ds.Empty))),
                      child: const ListTile(
                        leading: Icon(Icons.drive_folder_upload_outlined),
                        title: Text("Automatic Archival"),
                      ),
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
          meta.UploadsRow(margin: defaults.margin.copyWith(top: 0, bottom: 0) * 2),
        ],
      ),
      children: items,
      ds.Table.expanded<media.Media>(
        (v) {
          final onChange = (media.Media? upd) {
            setState(() {
              _res = api.FilesystemSearchResponse(
                next: _res.next,
                breadcrumb: _res.breadcrumb,
                items: ds.fnOnChange(_res.items, upd, (o) => o.id == v.id),
              );
            });
          };

          final preview = (BuildContext context, media.Media v) => media.Preview.modal(
            context,
            v,
            trailing: [
              ds.LoadingIconButton(
                icon: const Icon(Icons.delete_outline),
                onPressed: () => confirmremove(
                  context,
                  v,
                  apiremove: widget.apiremove,
                  onChange: (upd) {
                    // close the preview before the row is dropped from the listing.
                    ds.modals.push(context, null);
                    onChange(upd);
                  },
                ),
              ),
            ],
          );

          return media.RowDisplay(
            media: v,
            highlighted: v.id == _focused,
            leading: [Icon(mimex.icon(v.mimetype))],
            onTap: () {
              if (v.mimetype == mimex.directory) return Future.sync(() => navigate(v.id));
              final play = media.PlayAction(context, v, media.media.response());
              // nothing to play or preview: open the details instead.
              return Future.sync(play ?? preview(context, v) ?? focus(v));
            },
            trailing: [
              ds.Help(
                IconButton(
                  icon: const Icon(Icons.info_outline),
                  onPressed: focus(v),
                ),
                ds.Hint(const Text("show or hide details: type, timestamps, and actions like delete")),
              ),
            ],
            expanded: v.id == _focused ? FilesystemDetails(v, onChange: onChange) : ds.Empty,
          );
        },
      ),
    );

    return ds.Overlay(
      table,
      overlay: ds.Empty.maybe(
        _overlay,
        // an alignment expands the container to cover the browser completely.
        ds.Container(_overlay, alignment: Alignment.topCenter),
      ),
    );
  }

  void overlay(Widget w) {
    setState(() {
      _overlay = w;
    });
  }

  void mkdir() {
    ds.modals.push(
      context,
      DirectoryCreate(
        parent: directory,
        create: widget.apicreate,
        onCancel: () => ds.modals.push(context, null),
        onCreated: (created) {
          ds.modals.push(context, null);
          refresh(_res.next);
        },
      ),
    );
  }
}
