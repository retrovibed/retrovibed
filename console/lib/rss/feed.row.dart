import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/storage.dart' as storage;
import 'package:retrovibed/uuidx.dart' as uuidx;
import './api.dart' as api;

class FeedRow extends StatefulWidget {
  final api.Feed current;
  final Function(api.Feed?)? onChange;
  FeedRow({super.key, api.Feed? current, this.onChange})
    : current = current ?? (api.Feed.create()..autodownload = false);

  @override
  State<FeedRow> createState() => _FeedRowState();
}

class _FeedRowState extends State<FeedRow> {
  Widget _cause = ds.Error.zero;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _resetCause() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final placeholder = defaults.isCompact ? SizedBox(width: 0.0, height: 0.0) : SizedBox(width: 24.0, height: 24.0);

    return ds.ErrorScreen(
      ds.CompactingMenu(
        Text(
          widget.current.description,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
        ),
        icon: Icon(Icons.expand_more_rounded),
        trailing: [
          if (widget.current.hasNextCheck())
            ds.CompactingMenu.pinned(ds.Duration.untilISO8601(widget.current.nextCheck)),
          widget.current.autodownload ? Icon(Icons.downloading_rounded) : placeholder,
          widget.current.autoarchive ? Icon(Icons.archive_outlined) : placeholder,
          widget.current.contributing ? Icon(Icons.handshake_outlined) : placeholder,
          storage.SeedIcon(
            widget.current.encryptionSeed,
            classifier: storage.Classifier(
              community: widget.current.encryptionSeed,
              personal: uuidx.max(),
            ),
          ),
          ds.LoadingIconButton(
            icon: Icon(Icons.refresh),
            tooltip: "force refresh",
            onPressed: () {
              return httpx
                  .withRetry(() {
                    return api.refresh(
                      api.FeedCreateRequest(feed: widget.current),
                      options: [authn.AuthzCache.bearer(context)],
                    );
                  })
                  .then((resp) {
                    widget.onChange?.call(resp.feed);
                  })
                  .catchError((cause) {
                    setState(() {
                      _cause = ds.Error.unknown(cause, onTap: _resetCause);
                    });
                  });
            },
          ),
          ds.LoadingIconButton(
            icon: Icon(Icons.delete_outline),
            tooltip: "delete feed",
            onPressed: () async {
              ds.modals
                  .of(context)
                  ?.push(
                    ds.Confirmation.yesNo(
                      content: Text(
                        "Are you sure you want to delete the ${widget.current.description} feed?",
                      ),
                      onConfirm: () {
                        httpx
                            .withRetry(
                              () => api.delete(
                                widget.current.id,
                                options: [authn.AuthzCache.bearer(context)],
                              ),
                            )
                            .then((resp) => widget.onChange?.call(null))
                            .catchError(
                              (cause) => widget.onChange?.call(null),
                              test: httpx.ErrorsTest.err404,
                            )
                            .catchError((cause) {
                              setState(() {
                                _cause = ds.Error.unknown(
                                  cause,
                                  onTap: _resetCause,
                                );
                              });
                            })
                            .whenComplete(() {
                              ds.modals.of(context)?.push(null);
                            });
                      },
                      onCancel: () {
                        ds.modals.of(context)?.push(null);
                      },
                    ),
                  );
            },
          ),
        ],
      ),
      cause: _cause,
      tint: defaults.dangerTint,
      borderRadius: defaults.borderRadius,
    );
  }
}
