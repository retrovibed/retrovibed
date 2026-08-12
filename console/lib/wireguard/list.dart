import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/authn.dart' as authn;
import './meta.wireguard.pb.dart';
import './api.dart' as api;
import './list.row.dart';
import 'package:retrovibed/design.kit/stateful.dart';

class ListDisplay extends StatefulWidget {
  final api.FnWireguardSearch search;
  final api.FnUploadRequest upload;
  final TextEditingController? controller;
  final FocusNode? focus;
  const ListDisplay({
    super.key,
    this.search = api.wireguard.get,
    this.upload = api.wireguard.upload,
    this.controller,
    this.focus,
  });

  @override
  State<StatefulWidget> createState() => _ListDisplay();
}

class _ListDisplay extends State<ListDisplay> with LoadingState {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  Wireguard _current = Wireguard();
  api.WireguardSearchResponse _res = api.wireguard.response(
    next: api.wireguard.request(limit: 32),
  );

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> refresh(api.WireguardSearchRequest req) {
    return widget
        .search(req)
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
    _res.next..query = widget.controller?.text ?? "";
    refresh(_res.next);
    api.wireguard
        .current()
        .then(
          (r) => setState(() {
            _current = r.wireguard;
          }),
        )
        .catchError((cause) {}, test: httpx.ErrorsTest.err404)
        .catchError((cause) {
          print("failed to load current vpn settings ${cause}");
        })
        .ignore();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final upload =
        (
          ds.FilesEvent v, {
          ValueNotifier<int>? progress,
        }) {
          setState(() {
            _loading = true;
          });

          final multiparts = v.files.map((c) {
            return api.wireguard.uploadable(c.path, c.name, c.mimeType!);
          });

          return Future.microtask(() {
            return Future.wait(
                  multiparts.map((fv) {
                    return fv.then((v) {
                      return widget
                          .upload((req) {
                            req..files.add(v);
                            return req;
                          })
                          .then(
                            (v) => api.wireguard
                                .touch(
                                  v.wireguard.id,
                                  options: [authn.request(authn.AuthzCache.meta(context))],
                                )
                                .then((_) => v),
                          )
                          .then((uploaded) {
                            setState(() {
                              _res.items.add(uploaded.wireguard);
                              _current = uploaded.wireguard;
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
        controller: widget.controller,
        focus: widget.focus,
        onSubmitted: (v) {
          setState(() {
            _res.next.query = v;
            _res.next.offset = ds.SearchTray.Zero;
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
        empty: _res.items.length < _res.next.limit.toInt(),
        leading: [ds.FileDropWell.icon(upload, icon: Icons.add)],
        autofocus: defaults.desktop,
      ),
      children: _res.items,
      ds.Table.expanded<api.Wireguard>((v) {
        final onTap = () {
          return api.wireguard
              .touch(
                _current.id == v.id ? uuidx.max() : v.id,
                options: [authn.request(authn.AuthzCache.meta(context))],
              )
              .then((r) {
                setState(() {
                  _current = r.wireguard;
                });
              })
              .catchError((cause) {
                setState(() {
                  _current = Wireguard();
                });
              }, test: httpx.ErrorsTest.err404)
              .catchError((cause) {
                print("unexpected wireguard failure ${cause}");
                setState(() {
                  _current = Wireguard();
                });
              });
        };
        return ListRow(
          v,
          key: ValueKey(v.id),
          active: _current.id == v.id,
          onTap: onTap,
          onChange: (upd) {
            final updated = api.WireguardSearchResponse(
              items: ds.fnOnChange(_res.items, upd, (wg) => wg.id == upd.id),
              next: _res.next,
            );

            setState(() {
              _res = updated;
            });
          },
          onDelete: (deleted) {
            final updated = api.WireguardSearchResponse(
              items: ds.fnOnChange(_res.items, null, (wg) => wg.id == deleted.id),
              next: _res.next,
            );
            setState(() {
              _res = updated;
            });
          },
        );
      }),
      empty: ds.FileDropWell(
        upload,
        child: ds.FileDropWell.textual("drop a wireguard configuration file"),
      ),
    );
  }
}
