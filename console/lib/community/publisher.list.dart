import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import './api.dart' as api;
import './publisher.list.row.dart';

// The installed publish plugins. Unlike the search plugin list there is no
// search tray: the catalog endpoint returns every installed module in one
// shot, so paging controls would be inert.
class ListDisplay extends StatefulWidget {
  const ListDisplay({super.key});

  @override
  State<StatefulWidget> createState() => _ListDisplay();
}

class _ListDisplay extends State<ListDisplay> with ds.LoadingState {
  api.SocialsSearchResponse _res = api.SocialsSearchResponse();

  Future<void> refresh() {
    setState(() => loading = true);
    return api.publishers
        .search(options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _res = v;
            cause = ds.Error.zero;
          });
        })
        .catchError((cause) {}, test: httpx.ErrorsTest.err404)
        .catchError((cause) {
          setState(() => this.cause = ds.Error.unauthorized(cause, onTap: reseterr));
        }, test: httpx.ErrorsTest.unauthorized)
        .catchError((cause) {
          setState(() => this.cause = ds.Error.unknown(cause, onTap: reseterr));
        })
        .whenComplete(() => setState(() => loading = false));
  }

  @override
  void initState() {
    super.initState();
    ds.postframe(refresh);
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
            loading = true;
          });

          // .wasm rarely carries a mimetype from the platform's file picker,
          // and the endpoint rejects a blank one.
          final uploads = v.files.map((c) {
            return api.publishers
                .uploadable(c.path, c.name, c.mimeType ?? "application/wasm")
                .then((file) => (file, c.mimeType ?? "application/wasm"));
          });

          return Future.microtask(() {
            return Future.wait(
                  uploads.map((fv) {
                    return fv.then((v) {
                      final (file, mimetype) = v;
                      return api.publishers
                          .upload(
                            file.filename!.replaceAll(RegExp(r'\.wasm$'), ''),
                            mimetype,
                            (req) {
                              req.files.add(file);
                              return req;
                            },
                            options: [authn.request(authn.AuthzCache.meta(context))],
                          )
                          .then((created) {
                            setState(() {
                              _res.catalog.add(created.publisher);
                            });
                          })
                          .catchError((cause) {
                            setState(() {
                              this.cause = ds.Error.unknown(cause, onTap: reseterr);
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
      leading: Row(
        mainAxisAlignment: MainAxisAlignment.end,
        children: [
          ds.FileDropWell.icon(
            upload,
            icon: Icons.add,
            extensions: const ['wasm'],
            help: ds.Hint(const Text("install a publishing plugin")),
          ),
        ],
      ),
      children: _res.catalog,
      ds.Table.expanded<api.PluginPublisher>((v) {
        return PublisherRow(
          v,
          key: ValueKey(v.id),
          onDelete: (deleted) {
            setState(() {
              _res.catalog.removeWhere((p) => p.id == deleted.id);
            });
          },
        );
      }),
      empty: ds.FileDropWell(
        upload,
        extensions: const ['wasm'],
        child: ds.FileDropWell.textual("drop a compiled .wasm publisher"),
        shape: RoundedRectangleBorder(borderRadius: defaults.borderRadius),
      ),
    );
  }
}
