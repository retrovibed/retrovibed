import 'dart:async';

import 'package:retrovibed/design.kit/file.drop.well.dart';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/lucene.dart' as lucene;
import 'package:retrovibed/torrentx/display.dart' as torrentx;
import 'download.display.dart';
import 'grid.settings.dart';
import 'magnet.links.dart';

class AvailableListDisplay extends StatefulWidget {
  final media.FnDownloadSearch search;
  final media.FnUploadRequest upload;
  final TextEditingController? controller;
  final StreamController<media.Download> events;
  final List<Widget> leading;
  final List<Widget> trailing;
  const AvailableListDisplay(
    this.events, {
    super.key,
    this.search = media.discovered.available,
    this.upload = media.discovered.upload,
    this.controller,
    this.leading = const [],
    this.trailing = const [],
  });

  @override
  State<StatefulWidget> createState() => _AvailableListDisplay();
}

class _AvailableListDisplay extends State<AvailableListDisplay> with ds.LoadingState {
  StreamSubscription<void>? subscription;
  String _focused = '';
  Widget _tuning = ds.Empty;
  media.DownloadSearchResponse _res = media.discoveredsearch.response(
    next: media.discoveredsearch.request(limit: 32),
  );

  Future<void> refresh(media.DownloadSearchRequest req) {
    return widget
        .search(req, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _res = v;
            loading = false;
          });
        })
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
    subscription = widget.events.stream.listen((v) {
      // an empty download signals a full refresh instead of an in place update.
      if (v.media.id == "") {
        refresh(_res.next);
        return;
      }

      setState(() {
        _res = media.DownloadSearchResponse(
          items: ds.fnOnChange(_res.items, v, (d) => d.media.id == v.media.id),
          next: _res.next,
        );
      });
    });
    ds.postframe(() => refresh(_res.next));
  }

  @override
  void dispose() {
    subscription?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final upload =
        (
          FilesEvent v, {
          StreamSink<httpx.UploadProgress>? progress,
        }) {
          setState(() {
            loading = true;
          });

          return Future.microtask(() {
            final multiparts = v.files.map((c) {
              return media.media.uploadable(c.path, c.name, c.mimeType!);
            });
            return Future.wait(
                  multiparts.map((fv) {
                    return fv.then((v) {
                      return widget
                          .upload((req) {
                            req..files.add(v);
                            return req;
                          })
                          .then((uploaded) {
                            return media.discovered.download(uploaded.media.id);
                          });
                    });
                  }),
                )
                .then((v) => ds.NullWidget)
                .catchError((cause) {
                  return ds.Error.unknown(cause, onTap: reseterr);
                })
                .whenComplete(() => widget.events.add(media.Download()));
          });
        };

    return ds.Table(
      loading: loading,
      cause: cause,
      children: _res.items,
      leading: Column(
        verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
        children: [
          ds.SearchTray(
            autoscroll: true,
            autofocus: defaults.desktop,
            decoration: InputDecoration(hintText: "search downloadable content"),
            controller: widget.controller,
            filters: [
              lucene.Boolean.auto('completed', false, (v) {
                setState(() => _res.next.completed = v);
                refresh(_res.next);
              }),
              lucene.Boolean.auto('hidden', false, (v) {
                setState(() => _res.next.hidden = v);
                refresh(_res.next);
              }),
            ],
            onSubmitted: (v) {
              setState(() {
                _res.next.query = v;
                _res.next.offset = ds.Table.offset(0);
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
            empty: ds.Table.offset(_res.items.length) < _res.next.limit,
            leading: [
              ...widget.leading,
              ds.FileDropWell.icon(
                upload,
                mimetypes: [mimex.bittorrent],
                tooltip: "upload",
                help: ds.Hint(Text("upload torrent files to download")),
              ),
              ds.buttons.link(
                onPressed: () {
                  ds.modals.push(
                    context,
                    MagnetDownloads(
                      onSubmitted: (magents) {
                        final pending = magents.map(
                          (v) => media.discovered.magnet(
                            media.MagnetCreateRequest(uri: v),
                            options: [authn.request(authn.AuthzCache.meta(context))],
                          ),
                        );
                        return Future.wait(pending, eagerError: true).then((_) {
                          widget.events.add(media.Download());
                          ds.modals.of(context)?.reset();
                        });
                      },
                    ),
                  );
                },
                help: ds.Hint(Text("upload magnet urls to download")),
              ),
            ],
            tuning: ds.buttons.settings(
              onPressed: () => setState(() {
                _tuning = _tuning == ds.Empty
                    ? GridSettings(
                        _res.next,
                        onChange: (media.DownloadSearchRequest n) {
                          setState(() {
                            _res.next = n;
                          });
                        },
                      )
                    : ds.Empty;
              }),
              help: ds.Hint(Text("display advance settings")),
            ),
            help: ds.Hint(const Text("search discovered content, use @ to access advanced filtering")),
          ),
          ...widget.trailing,
        ],
      ),
      ds.Table.expanded<media.Download>(
        (v) {
          return ds.KeyPressAware.delete(
            onPress: () {
              return media.discovered
                  .reset(v.media.id, options: [authn.request(authn.AuthzCache.meta(context))])
                  .then((_) {
                    widget.events.add(v);
                  })
                  .catchError((cause) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text("Failed to reset: ${cause}")),
                    );
                    return null;
                  });
            },
            Column(
              mainAxisSize: MainAxisSize.max,
              children: [
                media.RowDisplay(
                  media: v.media,
                  leading: [Icon(mimex.icon(v.media.mimetype))],
                  help: ds.Hint.multiline([
                    Text("A downloadable media item."),
                    ds.HelpLabelled(
                      label: Text("tap"),
                      description: Text("expand details: file path, size, and distribution status"),
                    ),
                    ds.HelpLabelled(
                      label: Text("delete"),
                      description: Text("remove the item"),
                    ),
                  ]),
                  highlighted: _focused == v.media.id,
                  onTap: () async {
                    setState(() {
                      _focused = _focused == v.media.id ? '' : v.media.id;
                    });
                  },
                  trailing: [
                    ds.LoadingIconButton(
                      icon: Icon(media.download.icon(v)),
                      disabled: media.download.ongoing(v) && !media.download.paused(v),
                      help: ds.Hint(
                        Text(
                          media.download.typography(v),
                        ),
                      ),
                      onPressed: () => media.discovered
                          .download(
                            v.media.id,
                            options: [authn.request(authn.AuthzCache.meta(context))],
                          )
                          .then((d) {
                            widget.events.add(d.download);
                          })
                          .catchError((cause) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(
                                content: Text("Failed to download: $cause"),
                              ),
                            );
                            return null;
                          }),
                    ),
                  ],
                  expanded: Column(
                    spacing: defaults.spacing / 2,
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      DownloadDisplay(
                        v,
                        onVerify: (download) => ds.modals.asyncfn(
                          context,
                          (completion) => ds.Confirmation.yesNo(
                            content: Text(
                              "Are you sure you want to verify ${v.media.description}?",
                            ),
                            onConfirm: (context) {
                              media.discovered
                                  .update(
                                    v.media.torrentId,
                                    download..verifyAt = DateTime.now().toUtc().toIso8601String(),
                                    options: [authn.request(authn.AuthzCache.meta(context))],
                                  )
                                  .then((_) => completion.complete())
                                  .catchError((cause) {
                                    completion.completeError(cause);
                                  });
                            },
                            onCancel: (_) => completion.complete(),
                          ),
                        ),
                        onReset: () => ds.modals.asyncfn(
                          context,
                          (completion) => ds.Confirmation.yesNo(
                            content: Text(
                              "Are you sure you want to reset ${v.media.description}?",
                            ),
                            onConfirm: (context) {
                              httpx
                                  .withRetry(
                                    () => media.discovered.reset(
                                      v.media.id,
                                      options: [
                                        authn.request(authn.AuthzCache.meta(context)),
                                      ],
                                    ),
                                  )
                                  .then((__v) {
                                    setState(() {
                                      _res = media.DownloadSearchResponse(
                                        items: ds.fnOnChange(_res.items, null, (d) => d.media.id == v.media.id),
                                        next: _res.next,
                                      );
                                    });
                                    completion.complete();
                                  })
                                  .catchError((cause) {
                                    completion.completeError(cause);
                                  });
                            },
                            onCancel: (_) => completion.complete(),
                          ),
                        ),
                        onDelete: () => ds.modals.asyncfn(
                          context,
                          ds.Confirmation.dangerous(
                            content: Text(
                              "Are you sure you want to permanently delete ${v.media.description}?",
                            ),
                            onConfirm: (ctx) => httpx
                                .withRetry(
                                  () => media.discovered.delete(
                                    v.media.id,
                                    options: [authn.request(authn.AuthzCache.meta(ctx))],
                                  ),
                                )
                                .then((_) {
                                  setState(() {
                                    _res = media.DownloadSearchResponse(
                                      items: ds.fnOnChange(_res.items, null, (d) => d.media.id == v.media.id),
                                      next: _res.next,
                                    );
                                  });
                                }),
                          ),
                        ),
                      ),
                      torrentx.TorrentDisplay.fromID(
                        v.media.torrentId,
                      ),
                    ],
                  ),
                ),
              ],
            ),
          );
        },
        leading: [_tuning],
      ),
    );
  }
}
