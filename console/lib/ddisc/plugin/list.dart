import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import './api.dart' as api;
import './list.row.dart';
import 'package:retrovibed/design.kit/stateful.dart';

class ListDisplay extends StatefulWidget {
  const ListDisplay({super.key});

  @override
  State<StatefulWidget> createState() => _ListDisplay();
}

class _ListDisplay extends State<ListDisplay> with LoadingState {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  api.PluginSearchResponse _res = api.plugins.response(
    next: api.plugins.request(limit: 32),
  );

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> refresh(BuildContext context, api.PluginSearchRequest req) {
    return api.plugins
        .search(req, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _res = v;
            _loading = false;
          });
        })
        .catchError((cause) {
          setState(() {
            _loading = false;
          });
        }, test: httpx.ErrorsTest.err404)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unauthorized(cause, onTap: reseterr);
            _loading = false;
          });
        }, test: httpx.ErrorsTest.unauthorized)
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
    ds.postframe(() => refresh(context, _res.next));
  }

  @override
  Widget build(BuildContext context) {
    final upload =
        (
          ds.FilesEvent v, {
          ValueNotifier<int>? progress,
        }) {
          setState(() {
            _loading = true;
          });

          final uploads = v.files.map((c) {
            return api.plugins.uploadable(c.path, c.name, c.mimeType!);
          });

          return Future.microtask(() {
            return Future.wait(
                  uploads.map((fv) {
                    return fv.then((file) {
                      return api.plugins
                          .upload(
                            file.filename!.replaceAll(RegExp(r'\.wasm$'), ''),
                            (req) {
                              req.files.add(file);
                              return req;
                            },
                            options: [authn.request(authn.AuthzCache.meta(context))],
                          )
                          .then((created) {
                            setState(() {
                              _res.items.add(created.plugin);
                            });
                          })
                          .catchError((cause) {
                            setState(() {
                              _cause = ds.Error.unknown(cause, onTap: reseterr);
                            });
                          });
                    });
                  }),
                )
                .then((v) => ds.NullWidget)
                .catchError((cause) {
                  return ds.Error.unknown(cause, onTap: reseterr);
                })
                .whenComplete(
                  () => setState(() {
                    _loading = false;
                  }),
                );
          });
        };

    return ds.Table(
      loading: _loading,
      cause: _cause,
      leading: ds.SearchTray(
        onSubmitted: (v) => refresh(context, _res.next),
        next: (i) {
          setState(() {
            _res.next.offset = i;
          });
          refresh(context, _res.next);
        },
        current: _res.next.offset,
        empty: _res.items.length < _res.next.limit.toInt(),
        leading: [ds.FileDropWell.icon(upload, icon: Icons.add, extensions: const ['wasm'])],
      ),
      children: _res.items,
      ds.Table.expanded<api.Plugin>((v) {
        return ListRow(
          v,
          key: ValueKey(v.id),
          onDelete: (deleted) {
            setState(() {
              _res.items.removeWhere((p) => p.id == deleted.id);
            });
          },
        );
      }),
      empty: ds.FileDropWell(
        upload,
        extensions: const ['wasm'],
        child: ds.FileDropWell.textual("drop a compiled .wasm plugin"),
      ),
    );
  }
}
