import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/mimex.dart' as mimex;
import 'nettype.icon.dart';
import 'meta.wireguard.pb.dart';
import 'api.dart' as api;
import 'list.row.dart';

class ListDisplay extends StatefulWidget {
  final api.FnWireguardSearch apisearch;
  final api.FnUploadRequest apiupload;
  final api.FnWireguardCurrent apicurrent;
  final TextEditingController? controller;
  final FocusNode? focus;
  const ListDisplay({
    super.key,
    this.apisearch = api.wireguard.get,
    this.apiupload = api.wireguard.upload,
    this.apicurrent = api.wireguard.current,
    this.controller,
    this.focus,
  });

  @override
  State<StatefulWidget> createState() => _ListDisplay();
}

class _ListDisplay extends State<ListDisplay> with ds.LoadingState {
  Wireguard _distribution = Wireguard();
  api.WireguardSearchResponse _res = api.wireguard.response(
    next: api.wireguard.request(limit: 32),
  );

  Future<void> refresh(api.WireguardSearchRequest req) {
    return widget
        .apisearch(req)
        .then((v) {
          setState(() {
            _res = v;
            loading = false;
          });

          widget.focus?.requestFocus();
          ds.textediting.refocus(widget.controller);
        })
        .catchError((cause) {
          setState(() {
            loading = false;
          });
        }, test: httpx.ErrorsTest.err404)
        .catchError((cause) {
          setState(() {
            cause = ds.Error.unauthorized(cause, onTap: reseterr);
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

  @override
  void initState() {
    super.initState();
    _res.next..query = widget.controller?.text ?? "";
    refresh(_res.next);
    api.wireguard
        .current(api.WireguardNettype.DISTRIBUTION)
        .then(
          (r) => setState(() {
            _distribution = r.wireguard;
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
    const mimetypes = [mimex.text.plain];
    final defaults = ds.Defaults.of(context);
    final upload =
        (
          ds.FilesEvent v, {
          ValueNotifier<int>? progress,
        }) {
          setState(() {
            loading = true;
          });

          final multiparts = v.files.map((c) {
            return api.wireguard.uploadable(c.path, c.name, c.mimeType!);
          });

          return Future.microtask(() {
            return Future.wait(
                  multiparts.map((fv) {
                    return fv.then((v) {
                      return widget
                          .apiupload((req) {
                            req..files.add(v);
                            return req;
                          })
                          .then(
                            (v) => api.wireguard
                                .touch(
                                  v.wireguard.id,
                                  v.wireguard.nettype,
                                  options: [authn.request(authn.AuthzCache.meta(context))],
                                )
                                .then((_) => v),
                          )
                          .then((uploaded) {
                            setState(() {
                              _res.items.add(uploaded.wireguard);
                              _distribution = uploaded.wireguard;
                            });
                          })
                          .catchError((c) {
                            setState(() {
                              cause = ds.Error.unknown(c, onTap: reseterr);
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
                    loading = false;
                  }),
                );
          });
        };

    return ds.Table(
      loading: loading,
      cause: cause,
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
        leading: [
          ds.FileDropWell.icon(
            upload,
            mimetypes: mimetypes,
            icon: Icons.add,
          ),
        ],
        autofocus: defaults.desktop,
      ),
      children: _res.items,
      empty: ds.FileDropWell(
        upload,
        margin: defaults.margin,
        mimetypes: mimetypes,
        child: ds.FileDropWell.textual("drop a wireguard configuration file"),
        shape: RoundedRectangleBorder(borderRadius: defaults.borderRadius),
      ),
      ds.Table.expanded<api.Wireguard>((v) {
        final onNettype = (api.WireguardNettype nettype) {
          return api.wireguard
              .touch(
                _distribution.id == v.id ? uuidx.max() : v.id,
                nettype,
                options: [authn.request(authn.AuthzCache.meta(context))],
              )
              .then((r) {
                setState(() {
                  _distribution = r.wireguard;
                });
              })
              .catchError((cause) {
                setState(() {
                  _distribution = Wireguard();
                });
              }, test: httpx.ErrorsTest.err404)
              .catchError((cause) {
                print("unexpected wireguard failure ${cause}");
                setState(() {
                  _distribution = Wireguard();
                });
              });
        };
        return ListRow(
          v,
          key: ValueKey(v.id),
          onChange: (_, upd) async {
            print("DERP DERP ${upd}");
            final updated = api.WireguardSearchResponse(
              items: ds.fnOnChange(_res.items, upd, (wg) => wg.id == upd.id),
              next: _res.next,
            );

            setState(() {
              _res = updated;
            });
          },
          onDelete: (deleted) async {
            final updated = api.WireguardSearchResponse(
              items: ds.fnOnChange(_res.items, null, (wg) => wg.id == deleted.id),
              next: _res.next,
            );
            setState(() {
              _res = updated;
            });
          },
          leading: [
            NettypeIcon(v, onTap: onNettype),
          ],
        );
      }),
    );
  }
}
