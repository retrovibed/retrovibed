import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import './api.dart' as api;
import 'device.manual.dart';

class DaemonList extends StatefulWidget {
  final Future<api.Daemon> Function(api.Daemon v)? onSelect;
  final void Function(api.Daemon d)? onRemove;
  final Future<api.DaemonSearchResponse> Function(api.DaemonSearchRequest) search;
  // excludes daemons that are this device (api.daemons.isLocalDevice) from
  // the rendered rows.
  final bool remoteonly;
  // hides the add/remove/default-downloads-target controls, leaving just the
  // search box and tappable rows.
  final bool readonly;
  final InputDecoration decoration;

  const DaemonList({
    super.key,
    this.search = api.daemons.search,
    this.onSelect,
    this.onRemove,
    this.remoteonly = false,
    this.readonly = false,
    this.decoration = const InputDecoration(hintText: "search servers"),
  });

  @override
  State<StatefulWidget> createState() => _DaemonList();
}

class _DaemonList extends State<DaemonList> {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  Widget? _optional = null;
  api.DaemonSearchResponse _res = api.daemons.response();

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> refresh(api.DaemonSearchRequest req) {
    return widget
        .search(req)
        .then((v) {
          setState(() {
            _res = v;
            _loading = false;
          });
        })
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e, onTap: reseterr);
            _loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();
    refresh(_res.next);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final current = httpx.host();
    final replace = (api.Daemon v) {
      final replaced = _res.items.map((o) {
        final upd = o.id == v.id
            ? v
            : (o
                ..default_100 = false
                ..downloads = !v.downloads && o.downloads);
        return upd;
      });

      setState(() {
        _res = api.DaemonSearchResponse(items: replaced, next: _res.next);
      });

      return v;
    };

    final items = widget.remoteonly ? _res.items.where((d) => !api.daemons.isLocalDevice(d)).toList() : _res.items;

    return ds.Table(
      loading: _loading,
      cause: _cause,
      children: items,
      empty: const Text("no other devices found"),
      leading: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          ds.SearchTray(
            autofocus: defaults.desktop,
            decoration: widget.decoration,
            onSubmitted: (v) {
              setState(() {
                _res.next.query = v;
                _res.next.offset = ds.Int64(0);
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
              if (!widget.readonly)
                IconButton(
                  onPressed: () {
                    setState(() {
                      _optional = _optional != null
                          ? null
                          : ManualConfiguration(
                              (daemon) {
                                setState(() {
                                  _optional = null;
                                  _res.next.offset = ds.Int64(0);
                                });
                                refresh(_res.next);
                              },
                            );
                    });
                  },
                  icon: Icon(_optional == null ? Icons.add : Icons.remove),
                ),
            ],
          ),
          ?_optional,
        ],
      ),
      ds.Table.expanded<api.Daemon>(
        (v) => ds.ErrorBoundary(
          _RowDisplay(
            hostname: current,
            current: v,
            readonly: widget.readonly,
            replace: replace,
            onSelect: widget.onSelect == null
                ? null
                : () {
                    return widget.onSelect!(v).then(replace).whenComplete(() {
                      setState(() {});
                    });
                  },
            onRemove: widget.readonly
                ? null
                : (api.Daemon d) {
                    return api.daemons.delete(d.id).then((v) {
                      _res.next.offset = _res.next.offset - 1;
                      refresh(_res.next);
                      return v.daemon;
                    });
                  },
          ),
        ),
      ),
    );
  }
}

class _RowDisplay extends StatelessWidget {
  final String hostname;
  final api.Daemon current;
  final bool readonly;
  final Future<api.Daemon> Function()? onSelect;
  final Future<api.Daemon> Function(api.Daemon d)? onRemove;
  final Function(api.Daemon d) replace;
  const _RowDisplay({
    required this.current,
    required this.hostname,
    required this.replace,
    this.readonly = false,
    this.onSelect,
    this.onRemove,
  });

  @override
  Widget build(BuildContext context) {
    return ds.TableRow(
      onTap: onSelect == null
          ? null
          : () {
              onSelect!()
                  .catchError(
                    ds.Error.boundary(context, current, ds.Error.offline),
                    test: ds.ErrorTests.offline,
                  )
                  .catchError(
                    ds.Error.boundary(
                      context,
                      current,
                      ds.Error.connectivity,
                    ),
                    test: ds.ErrorTests.connectivity,
                  )
                  .catchError(
                    ds.Error.boundary(context, current, ds.Error.unknown),
                  );
            },
      [
        Opacity(
          opacity: current.hostname == hostname || current.default_100 ? 1.0 : 0.0,
          child: Icon(Icons.check, color: Colors.lightGreenAccent),
        ),
        if (!readonly)
          ds.LoadingIconButton(
            disabled: current.downloads,
            onPressed: () {
              return api.daemons
                  .download(current.id, current)
                  .then((v) {
                    replace(v.daemon);
                    return v;
                  })
                  .catchError((
                    cause,
                  ) {
                    ds.ErrorBoundary.of(
                      context,
                    )?.onError(ds.Error.unknown(cause));
                    throw cause;
                  });
            },
            icon: Icon(Icons.downloading_rounded, color: current.downloads ? Colors.lightGreenAccent : null),
            tooltip: current.downloads
                ? "downloads will be sent to this library by default for background processes"
                : "mark this library the default for downloads from background processes",
          ),
        Expanded(
          child: Text(current.description, overflow: TextOverflow.ellipsis),
        ),
        if (onRemove != null)
          IconButton(
            onPressed: () {
              onRemove!(current).catchError((cause) {
                ds.ErrorBoundary.of(
                  context,
                )?.onError(ds.Error.unknown(cause));
                throw cause;
              }).ignore();
            },
            icon: Icon(Icons.remove),
          ),
      ],
    );
  }
}
