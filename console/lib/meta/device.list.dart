import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import './api.dart' as api;
import 'device.manual.dart';

class DaemonList extends StatefulWidget {
  final Future<api.Daemon> Function(api.Daemon v)? onDoubleTap;
  final void Function(api.Daemon d)? onRemove;
  final Future<api.DaemonSearchResponse> Function(api.DaemonSearchRequest) search;

  const DaemonList({
    super.key,
    this.search = api.daemons.search,
    this.onDoubleTap,
    this.onRemove,
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

    return ds.Table(
      loading: _loading,
      cause: _cause,
      children: _res.items,
      leading: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          ds.SearchTray(
            autofocus: defaults.desktop,
            decoration: InputDecoration(hintText: "search servers"),
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
      ds.Table.inline<api.Daemon>(
        (v) => ds.ErrorBoundary(
          _RowDisplay(
            hostname: current,
            current: v,
            replace: replace,
            onDoubleTap: widget.onDoubleTap == null
                ? null
                : () {
                    return api.daemons
                        .reachable(v)
                        .then((v) {
                          return widget.onDoubleTap!(v).then(replace);
                        })
                        .whenComplete(() {
                          setState(() {});
                        });
                  },
            onRemove: (api.Daemon d) {
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
  final Future<api.Daemon> Function()? onDoubleTap;
  final Future<api.Daemon> Function(api.Daemon d)? onRemove;
  final Function(api.Daemon d) replace;
  const _RowDisplay({
    required this.current,
    required this.hostname,
    required this.replace,
    this.onDoubleTap,
    this.onRemove,
  });

  @override
  Widget build(BuildContext context) {
    final themex = ds.Defaults.of(context);
    return Container(
      child: InkWell(
        onDoubleTap: onDoubleTap == null
            ? null
            : () {
                onDoubleTap!()
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
        child: Row(
          spacing: themex.spacing,
          children: [
            Opacity(
              opacity: current.hostname == hostname || current.default_100 ? 1.0 : 0.0,
              child: Icon(Icons.check, color: Colors.lightGreenAccent),
            ),
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
        ),
      ),
    );
  }
}
