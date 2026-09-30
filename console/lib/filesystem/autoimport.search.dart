import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
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
  final ValueNotifier<api.AutoimportDirectorySearchRequest> _next = ValueNotifier(api.autoimport.request());
  // number of results on the current page, the tray uses it to decide if there is a next page.
  int _count = 0;
  bool _creating = false;

  @override
  void dispose() {
    _next.dispose();
    super.dispose();
  }

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
            _next.value = _next.value.deepCopy()
              ..query = v
              ..offset = ds.Int64(0);
          },
          next: (i) {
            _next.value = _next.value.deepCopy()..offset = i;
          },
          current: _next.value.offset,
          empty: ds.Int64(_count) < _next.value.limit,
          leading: [
            ds.CompactingMenu.pinned(
              ds.LoadingIconButton.close(
                onPressed: () {
                  widget.onClose();
                  return Future.value(null);
                },
                help: ds.Hint(const Text("close monitored directories")),
              ),
            ),
            ds.Help(
              IconButton(
                onPressed: () => setState(() => _creating = !_creating),
                icon: Icon(_creating ? Icons.remove : Icons.add),
              ),
              ds.Hint(const Text("monitor a new directory")),
            ),
          ],
        ),
        Expanded(
          child: authn.AuthzCache(
            _AutoimportListing(
              next: _next,
              creating: _creating,
              onCreating: (v) => setState(() => _creating = v),
              onCount: (n) => setState(() => _count = n),
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
class _AutoimportListing extends StatefulWidget {
  final ValueNotifier<api.AutoimportDirectorySearchRequest> next;
  final bool creating;
  final void Function(bool) onCreating;
  final void Function(int) onCount;
  final api.FnAutoimportSearch search;
  final api.FnAutoimportCreate create;
  final api.FnAutoimportUpdate update;
  final api.FnAutoimportDelete delete;

  const _AutoimportListing({
    required this.next,
    required this.creating,
    required this.onCreating,
    required this.onCount,
    required this.search,
    required this.create,
    required this.update,
    required this.delete,
  });

  @override
  State<_AutoimportListing> createState() => _AutoimportListingState();
}

class _AutoimportListingState extends State<_AutoimportListing> with ds.LoadingState {
  api.AutoimportDirectory _created = api.autoimport.directory();
  List<api.AutoimportDirectory> _items = [];

  Future<void> refresh() {
    setState(() => loading = true);
    return widget
        .search(widget.next.value, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((r) {
          setState(() => _items = r.items);
          widget.onCount(r.items.length);
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
    widget.next.addListener(refresh);
    ds.postframe(refresh);
  }

  @override
  void dispose() {
    widget.next.removeListener(refresh);
    super.dispose();
  }

  void resetcreate() {
    setState(() => _created = api.autoimport.directory());
    widget.onCreating(false);
  }

  Future<void> submit(api.AutoimportDirectory n) {
    setState(() => loading = true);
    return widget
        .create(
          api.AutoimportDirectoryCreateRequest(directory: n),
          options: [authn.request(authn.AuthzCache.meta(context))],
        )
        .then((v) => resetcreate())
        .then((v) => refresh())
        .catchError((e) {
          setState(() {
            cause = ds.Error.unknown(e, onTap: reseterr);
            loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final proto = _AutoimportCreate(
      current: _created,
      onCancel: resetcreate,
      onSubmit: submit,
      onChange: (v) => setState(() => _created = v),
    );

    return ds.Table(
      loading: loading,
      cause: cause,
      padding: defaults.padding / 2,
      children: _items,
      ds.Table.expanded<api.AutoimportDirectory>(
        (d) => AutoimportItem(
          key: ValueKey(d.id),
          current: d,
          update: widget.update,
          delete: widget.delete,
          onChange: (v) {
            setState(() {
              _items = ds.fnOnChange(_items, v, (old) => old.id == d.id);
            });
          },
        ),
      ),
      empty: proto,
      overlay: widget.creating ? proto : ds.Empty,
    );
  }
}

class _AutoimportCreate extends StatelessWidget {
  final api.AutoimportDirectory current;
  final Function(api.AutoimportDirectory)? onChange;
  final Future<void> Function(api.AutoimportDirectory)? onSubmit;
  final Function()? onCancel;

  const _AutoimportCreate({required this.current, this.onChange, this.onCancel, this.onSubmit});

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
            AutoimportEdit(current: current, onChange: onChange, pathEditable: true),
            Row(
              spacing: defaults.spacing,
              children: [
                const Spacer(),
                ds.LoadingButton(const Text("cancel"), onPressed: () async => onCancel?.call()),
                ds.LoadingButton(const Text("create"), onPressed: () async => onSubmit?.call(current)),
                const Spacer(),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
