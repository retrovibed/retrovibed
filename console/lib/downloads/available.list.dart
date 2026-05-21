import 'package:retrovibed/design.kit/file.drop.well.dart';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/lucene.dart' as lucene;
import 'package:retrovibed/timex.dart' as timex;
import './grid.settings.dart';
import './magnet.links.dart';

class AvailableListDisplay extends StatefulWidget {
  final media.FnDownloadSearch search;
  final media.FnUploadRequest upload;
  final TextEditingController? controller;
  final ValueNotifier<int>? events;
  const AvailableListDisplay({
    super.key,
    this.search = media.discovered.available,
    this.upload = media.discovered.upload,
    this.controller,
    this.events,
  });

  @override
  State<StatefulWidget> createState() => _AvailableListDisplay();
}

class _AvailableListDisplay extends State<AvailableListDisplay> {
  bool _loading = true;
  String _focused = '';
  Widget _cause = ds.Error.zero;
  media.DownloadSearchResponse _res = media.discoveredsearch.response(
    next: media.discoveredsearch.request(limit: 32),
  );

  @override
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void resetcause() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> refresh(media.DownloadSearchRequest req) {
    return widget
        .search(req, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _res = v;
            _loading = false;
          });
        })
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e, onTap: resetcause);
            _loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();
    refresh(_res.next);
    widget.events?.addListener(() {
      refresh(_res.next);
    });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final theme = Theme.of(context);
    final upload = (
      FilesEvent v, {
      ValueNotifier<int>? progress,
    }) {
      setState(() {
        _loading = true;
      });

      return Future.microtask(() {
        final multiparts = v.files.map((c) {
          return media.media.uploadable(c.path, c.name, c.mimeType!);
        });
        return Future.wait(
              multiparts.map((fv) {
                return fv
                    .then((v) {
                      return widget
                          .upload((req) {
                            req..files.add(v);
                            return req;
                          })
                          .then((uploaded) {
                            return media.discovered.download(uploaded.media.id);
                          });
                    })
                    .whenComplete(() => widget.events?.value += 1);
              }),
            )
            .then((v) => ds.NullWidget)
            .catchError((cause) {
              return ds.Error.unknown(cause, onTap: resetcause);
            })
            .whenComplete(() {
              setState(() {
                _loading = false;
                widget.events?.value += 1;
              });
            });
      });
    };

    return ds.Table(
      loading: _loading,
      cause: _cause,
      children: _res.items,
      leading: ds.SearchTray(
        ensureVisible: true,
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
                      widget.events?.value += 1;
                      ds.modals.of(context)?.reset();
                    });
                  },
                ),
              );
            },
            help: ds.Hint(Text("upload magnet urls to download")),
          ),
        ],
        tuning: GridSettings(
          _res.next,
          onChange: (media.DownloadSearchRequest n) {
            setState(() {
              _res.next = n;
            });
          },
        ),
        help: ds.Hint(const Text("search discovered content, use @ to access advanced filtering")),
      ),
      ds.Table.expanded<media.Download>(
        (v) {
          final downloading = timex.iso8601(v.initiatedAt).isBefore(timex.inf);
          final paused = timex.iso8601(v.pausedAt).isBefore(timex.inf);

          return ds.KeyPressAware.delete(
            onPress: () {
              return media.discovered
                  .reset(v.media.id, options: [authn.request(authn.AuthzCache.meta(context))])
                  .then((v) {
                    widget.events ?? refresh(_res.next);
                    widget.events?.value += 1;
                  })
                  .catchError((cause) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text("Failed to delete: ${cause}")),
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
                  onTap: () async {
                    setState(() {
                      _focused = _focused == v.media.id ? '' : v.media.id;
                    });
                  },
                  trailing: [
                    ds.LoadingIconButton(
                      icon: Icon(downloading ? Icons.downloading : Icons.download),
                      disabled: downloading && !paused,
                      help: ds.Hint(
                        Text(
                          downloading && !paused
                              ? "Download is in progress."
                              : paused
                              ? "Resume the download for this item."
                              : "Start downloading this item.",
                        ),
                      ),
                      onPressed:
                          () => media.discovered
                              .download(
                                v.media.id,
                                options: [authn.request(authn.AuthzCache.meta(context))],
                              )
                              .then((v) {
                                widget.events ?? refresh(_res.next);
                                widget.events?.value += 1;
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
                ),
                if (_focused == v.media.id)
                  media.DownloadDisplay(
                    v,
                    background: theme.colorScheme.surfaceContainerLow,
                    onVerify:
                        (download) => ds.modals.asyncfn(
                          context,
                          (completion) => ds.Confirmation.yesNo(
                            content: Text(
                              "Are you sure you want to verify ${v.media.description}?",
                            ),
                            onConfirm: () {
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
                            onCancel: completion.complete,
                          ),
                        ),
                    onTap:
                        () => ds.modals.asyncfn(
                          context,
                          (completion) => ds.Confirmation.yesNo(
                            content: Text(
                              "Are you sure you want to reset ${v.media.description}?",
                            ),
                            onConfirm: () {
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
                                        items: _res.items.where(
                                          (d) => d.media.id != v.media.id,
                                        ),
                                        next: _res.next,
                                      );
                                    });
                                    completion.complete();
                                  })
                                  .catchError((cause) {
                                    completion.completeError(cause);
                                  });
                            },
                            onCancel: completion.complete,
                          ),
                        ),
                    trailing: [],
                  ),
              ],
            ),
          );
        },
      ),
    );
  }
}
