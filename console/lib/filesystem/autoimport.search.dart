import 'dart:async';

import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta.dart' as meta;
import 'autoimport.api.dart' as api;
import 'autoimport.edit.dart';
import 'autoimport.list.dart';

// monitored directories are paths on this machine, so authorization is always obtained from the
// local daemon regardless of which library is currently selected.
Future<meta.AuthzResponse> localauthz({String? host}) => meta.authz.current(host: httpx.localhost());

// searchable listing of the directories the local daemon monitors for automatic import.
//
// the search tray lives outside the local authorization cache: when the local daemon can't be
// reached the cache covers its subtree with the error, and the close button must stay usable.
class AutoimportSearch extends StatefulWidget {
  final VoidCallback onClose;
  final authn.FnAuthzCurrent authz;
  final api.FnAutoimportSearch search;
  final api.FnAutoimportCreate create;
  final api.FnAutoimportUpdate update;
  final api.FnAutoimportDelete delete;

  const AutoimportSearch({
    super.key,
    required this.onClose,
    this.authz = localauthz,
    this.search = api.autoimport.search,
    this.create = api.autoimport.create,
    this.update = api.autoimport.update,
    this.delete = api.autoimport.delete,
  });

  @override
  State<AutoimportSearch> createState() => _AutoimportSearchState();
}

class _AutoimportSearchState extends State<AutoimportSearch> {
  api.AutoimportDirectorySearchResponse _res = api.autoimport.response();
  Widget _overlay = ds.Empty;

  void resetoverlay() => setState(() => _overlay = ds.Empty);

  // a new next instance is what triggers the listing to search.
  void search(api.AutoimportDirectorySearchRequest next) => setState(() {
    _res = api.AutoimportDirectorySearchResponse(next: next, items: _res.items);
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return Column(
      mainAxisSize: MainAxisSize.max,
      children: [
        ds.SearchTray(
          autofocus: defaults.desktop,
          decoration: const InputDecoration(hintText: "search monitored directories"),
          onSubmitted: (v) async {
            search(
              _res.next.deepCopy()
                ..query = v
                ..offset = ds.Int64(0),
            );
          },
          next: (i) => search(_res.next.deepCopy()..offset = i),
          current: _res.next.offset,
          empty: ds.Int64(_res.items.length) < _res.next.limit,
          tuning: ds.CompactingMenu.pinned(
            ds.LoadingIconButton.close(
              onPressed: () {
                widget.onClose();
                return Future.value(null);
              },
              help: ds.Hint(const Text("close monitored directories")),
            ),
          ),
          leading: [
            ds.FileDropWell.icon(
              (evt, {progress}) async {
                setState(() {
                  final path = evt.files.first.path;
                  _overlay = _AutoImportCreate(
                    key: ValueKey(evt.files.first.path),
                    current: api.autoimport.directory()
                      ..path = path
                      ..description = path,
                    create: widget.create,
                    onCancel: resetoverlay,
                    onCreated: (_) {
                      resetoverlay();
                      search(_res.next.deepCopy());
                    },
                  );
                });
                return ds.NullWidget;
              },
              mimetypes: mimex.folders,
              icon: Icons.create_new_folder_outlined,
              help: ds.Hint(const Text("select or drop a directory to monitor")),
            ),
          ],
        ),
        Expanded(
          child: authn.AuthzCache(
            _AutoImportListing(
              padding: defaults.padding / 2,
              margin: defaults.margin / 2,
              current: _res,
              overlay: _overlay,
              onChange: (r) => setState(() => _res = r),
              search: widget.search,
              create: widget.create,
              update: widget.update,
              delete: widget.delete,
            ),
            current: widget.authz,
          ),
        ),
      ],
    );
  }
}

// fetches and renders the monitored directories for the request published by the search tray.
// it lives below the local authorization cache so every request carries a local daemon token.
class _AutoImportListing extends StatefulWidget {
  final api.AutoimportDirectorySearchResponse current;
  final void Function(api.AutoimportDirectorySearchResponse) onChange;
  final api.FnAutoimportSearch search;
  final api.FnAutoimportCreate create;
  final api.FnAutoimportUpdate update;
  final api.FnAutoimportDelete delete;
  final Widget overlay;
  final EdgeInsets padding;
  final EdgeInsets margin;

  const _AutoImportListing({
    required this.current,
    this.overlay = ds.Empty,
    this.padding = EdgeInsets.zero,
    this.margin = EdgeInsets.zero,
    this.onChange = ds.fnNoop,
    required this.search,
    required this.create,
    required this.update,
    required this.delete,
  });

  @override
  State<_AutoImportListing> createState() => _AutoImportListingState();
}

class _AutoImportListingState extends State<_AutoImportListing> with ds.LoadingState {
  Future<void> refresh() {
    final next = widget.current.next;
    setState(() => loading = true);
    return widget
        .search(next, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((r) {
          // keep the same next instance so reporting the results doesn't trigger another search.
          widget.onChange(api.AutoimportDirectorySearchResponse(next: next, items: r.items));
        })
        .catchError((e) {
          setState(() {
            cause = ds.Error.unknown(e, onTap: reseterr);
          });
        })
        .whenComplete(() {
          setState(() => loading = false);
        });
  }

  @override
  void initState() {
    super.initState();
    ds.postframe(refresh);
  }

  @override
  void didUpdateWidget(covariant _AutoImportListing old) {
    super.didUpdateWidget(old);
    if (!identical(old.current.next, widget.current.next)) ds.postframe(refresh);
  }

  // monitors each dropped directory with the default debounce and mode. dropped paths are local,
  // matching the local daemon every autoimport request is sent to.
  Future<Widget?> dropped(ds.FilesEvent evt, {StreamSink<httpx.UploadProgress>? progress}) {
    final authz = authn.request(authn.AuthzCache.meta(context));
    return evt.files
        .fold<Future<void>>(Future.value(), (prev, f) {
          return prev.then(
            (_) => widget.create(
              api.AutoimportDirectoryCreateRequest(directory: api.autoimport.directory()..path = f.path),
              options: [authz],
            ),
          );
        })
        .then((_) => refresh())
        .then((_) => ds.NullWidget)
        .catchError((e) {
          setState(() {
            cause = ds.Error.unknown(e, onTap: reseterr);
          });
          return ds.NullWidget;
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return ds.Table(
      loading: loading,
      cause: cause,
      margin: widget.margin,
      padding: widget.padding,
      children: widget.current.items,
      ds.Table.expanded<api.AutoimportDirectory>(
        (d) => AutoImportItem(
          key: ValueKey(d.id),
          current: d,
          update: widget.update,
          delete: widget.delete,
          onChange: (v) => widget.onChange(
            api.AutoimportDirectorySearchResponse(
              next: widget.current.next,
              items: ds.fnOnChange(widget.current.items, v, (old) => old.id == d.id),
            ),
          ),
        ),
      ),
      empty: ds.FileDropWell(
        dropped,
        shape: RoundedRectangleBorder(borderRadius: defaults.borderRadius),
        mimetypes: mimex.folders,
        child: const Center(
          child: Column(
            mainAxisSize: MainAxisSize.max,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(mimex.icofolder),
              SelectableText("Drop a directory to monitor it for automatic import."),
            ],
          ),
        ),
      ),
      overlay: widget.overlay,
    );
  }
}

// rendered as the listing's overlay, so it sits below the local authorization cache.
class _AutoImportCreate extends StatefulWidget {
  final api.AutoimportDirectory current;
  final api.FnAutoimportCreate create;
  final VoidCallback onCancel;
  final void Function(api.AutoimportDirectory) onCreated;

  const _AutoImportCreate({
    super.key,
    required this.current,
    required this.create,
    required this.onCancel,
    required this.onCreated,
  });

  @override
  State<_AutoImportCreate> createState() => _AutoImportCreateState();
}

class _AutoImportCreateState extends State<_AutoImportCreate> with ds.LoadingState {
  late api.AutoimportDirectory _current = widget.current;

  Future<void> submit() {
    return widget
        .create(
          api.AutoimportDirectoryCreateRequest(directory: _current),
          options: [authn.request(authn.AuthzCache.meta(context))],
        )
        .then((v) => widget.onCreated(v.directory))
        .catchError((e) {
          setState(() {
            cause = ds.Error.unknown(e, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ds.Card(
      padding: EdgeInsets.zero,
      forms.Container(
        padding: defaults.padding,
        Column(
          spacing: defaults.spacing,
          mainAxisSize: MainAxisSize.min,
          children: [
            ds.ErrorScreen(
              cause: cause,
              AutoImportEdit(current: _current, onChange: (v) => setState(() => _current = v), pathEditable: true),
            ),
            Row(
              spacing: defaults.spacing,
              children: [
                const Spacer(),
                ds.LoadingButton(const Text("cancel"), onPressed: () async => widget.onCancel()),
                ds.LoadingButton(const Text("create"), onPressed: submit),
                const Spacer(),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
