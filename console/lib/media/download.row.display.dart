import 'dart:async';
import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import './api.dart' as api;
import './download.row.controls.dart';

class RefreshingDownload extends StatefulWidget {
  final api.Download current;
  final Duration interval;
  final api.FnDownloadWatch watch;
  const RefreshingDownload({
    super.key,
    required this.current,
    this.interval = const Duration(milliseconds: 5000),
    this.watch = api.discovered.watch,
  });

  @override
  State<RefreshingDownload> createState() => _DownloadingState();
}

class _DownloadingState extends State<RefreshingDownload> {
  Widget _cause = ds.Error.zero;
  api.Download current = api.Download();
  StreamSubscription<api.Download>? _subscription;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _resetcause() {
    setState(() => _cause = ds.Error.zero);
  }

  @override
  void initState() {
    super.initState();
    current = widget.current;
    widget
        .watch(
          widget.current.media.id,
          options: [authn.AuthzCache.bearer(context)],
        )
        .then((socket) {
          final c = Completer();

          _subscription = socket.listen(
            (v) {
              setState(() {
                current = v;
              });
            },
            cancelOnError: true,
            onError: c.completeError,
            onDone: c.complete,
          );

          return c.future;
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: _resetcause);
          });
        });
  }

  @override
  void dispose() {
    _subscription?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ds.ErrorScreen(
      cause: _cause,
      DownloadRowDisplay(
        current: current,
        trailing:
            (ctx) => DownloadRowControls(
              current: current,
              onChange: (d) {
                ds.RefreshBoundary.of(ctx)?.reset();
              },
            ),
      ),
    );
  }
}

class DownloadRowDisplay extends StatelessWidget {
  final api.Download current;
  final Widget? Function(BuildContext)? trailing;
  const DownloadRowDisplay({super.key, required this.current, this.trailing});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    var percentage = math.min(
      current.bytes == 0 ? 0.0 : (current.downloaded.toDouble() / current.bytes.toDouble()),
      1.0,
    );

    return ds.ErrorBoundary(
      SelectionArea(
        child: Builder(
          builder: (context) {
            final compact = defaults.isCompact;
            if (compact) {
              return Column(
                children: [
                  Row(
                    spacing: defaults.spacing,
                    children: [
                      Expanded(
                        child: Text(
                          current.media.description,
                          overflow: TextOverflow.ellipsis,
                          maxLines: 1,
                        ),
                      ),
                      Expanded(
                        child: LinearProgressIndicator(
                          value: percentage,
                          semanticsLabel: 'Linear progress indicator',
                        ),
                      ),
                    ],
                  ),
                  Row(
                    spacing: defaults.spacing,
                    children: [
                      Expanded(child: Icon(Icons.people_outline, size: 16)),
                      Expanded(child: Text(current.peers.toString())),
                      Expanded(child: Text("${(percentage * 100).toStringAsFixed(1)}%")),
                      trailing?.call(context) ?? const SizedBox(),
                    ],
                  ),
                ],
              );
            }

            return Row(
              spacing: defaults.spacing,
              children: [
                Expanded(
                  child: Text(
                    current.media.description,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                Expanded(
                  child: LinearProgressIndicator(
                    value: percentage,
                    semanticsLabel: 'Linear progress indicator',
                  ),
                ),
                Icon(Icons.people_outline, size: 16),
                Text(current.peers.toString()),
                Text("${(percentage * 100).toStringAsFixed(1)}%"),
                trailing?.call(context) ?? const SizedBox(),
              ],
            );
          },
        ),
      ),
    );
  }
}
