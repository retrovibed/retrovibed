import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import '../dropdown.upload.dart';
import '../list.display.dart';
import 'folder.create.dart';
import 'row.dart';

// navigates the library as a tree. AvailableListDisplay owns its own request and exposes
// no way to seed a parent, so the folder being listed is injected through the search
// closure and a change of folder is published by rekeying the list, which is what makes it
// refetch. everything else reuses that widget untouched so this view carries exactly the
// one row of chrome the library view does.
class FilesystemBrowser extends StatefulWidget {
  final media.FnMediaSearch search;
  final media.FnUploadRequest upload;
  final media.FnMkdir mkdir;
  final TextEditingController? controller;
  final FocusNode? focus;
  final ValueNotifier<media.SearchMode> mode;
  final void Function(media.SearchMode) onModeChanged;

  const FilesystemBrowser({
    super.key,
    required this.mode,
    required this.onModeChanged,
    this.search = media.media.search,
    this.upload = media.media.upload,
    this.mkdir = media.media.mkdir,
    this.controller,
    this.focus,
  });

  @override
  State<StatefulWidget> createState() => _FilesystemBrowser();
}

class _FilesystemBrowser extends State<FilesystemBrowser> {
  String _parent = uuidx.min();
  // creating a folder does not change the folder being listed, so the listing is rekeyed
  // explicitly to pick it up.
  int _generation = 0;
  String _highlighted = "";
  media.MediaSearchResponse _last = media.media.response();

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void navigate(String id) {
    setState(() {
      _parent = id;
      _highlighted = "";
    });
  }

  // the folder holding the one being listed. the daemon returns the path root first, so
  // the entry before the last is the parent; a single entry means the parent is the root.
  String get ancestor {
    final path = _last.breadcrumb;
    return path.length < 2 ? uuidx.min() : path[path.length - 2].id;
  }

  String get location {
    final path = _last.breadcrumb.map((v) => v.description).join("/");
    return path.isEmpty ? "search the library" : "search in ${path}";
  }

  Future<media.MediaSearchResponse> search(
    media.MediaSearchRequest req, {
    String? host,
    List<httpx.Option> options = const [],
  }) {
    req.parentId = _parent;

    return widget.search(req, host: host, options: options).then((v) {
      setState(() {
        _last = v;
      });

      // a parent entry leads the listing instead of a breadcrumb bar, so moving up costs
      // no chrome and renders through the same row widget as everything else.
      if (uuidx.fromString(v.next.parentId) != uuidx.fromString(uuidx.min())) {
        v.items.insert(0, media.Media(id: ancestor, description: "..", mimetype: mimex.directory));
      }

      return v;
    });
  }

  Widget row(media.Media v) {
    if (v.mimetype == mimex.directory) {
      return FilesystemRow(
        current: v,
        highlighted: v.id == _highlighted,
        onTap: () => Future.sync(() => navigate(v.id)),
      );
    }

    return FilesystemRow(
      current: v,
      highlighted: v.id == _highlighted,
      trailing: [media.ButtonShare(current: v)],
      // playback owns audio and video; everything else the reader can inspect without
      // pulling the whole file down first.
      onTap: media.PlayAction(context, v, _last) ?? preview(v),
    );
  }

  Future<void> Function() preview(media.Media v) {
    return () => Future.sync(
      () => ds.modals.push(
        context,
        ds.Card(
          leading: [ds.Heading(Text(v.description))],
          trailing: [
            IconButton(
              icon: const Icon(Icons.close),
              onPressed: () => ds.modals.push(context, null),
            ),
          ],
          SingleChildScrollView(child: media.Preview(current: v)),
        ),
      ),
    );
  }

  void mkdir() {
    ds.modals.push(
      context,
      FolderCreate(
        parent: _parent,
        mkdir: widget.mkdir,
        onCancel: () => ds.modals.push(context, null),
        onCreated: (created) {
          ds.modals.push(context, null);
          setState(() {
            _highlighted = created.id;
            _generation++;
          });
        },
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return AvailableListDisplay(
      key: ValueKey("$_parent/$_generation"),
      controller: widget.controller,
      focus: widget.focus,
      search: search,
      decoration: InputDecoration(hintText: location),
      // files dropped onto the listing belong to the folder on screen.
      upload: (mkreq) => widget.upload((req) {
        req.fields["parent_id"] = _parent;
        return mkreq(req);
      }),
      row: row,
      leading: [
        ds.CompactingMenu.pinned(
          DropdownUpload(
            icon: Icon(mimex.icofolder),
            help: ds.Hint(const Text("create a folder, or switch to library or discover mode")),
            items: [
              PopupMenuItem<String>(
                onTap: mkdir,
                child: ListTile(leading: Icon(mimex.icofolder), title: const Text("New Folder")),
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
      ],
    );
  }
}
