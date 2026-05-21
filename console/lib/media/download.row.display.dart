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

  void _reconnect() {
    if (!mounted) return;
    Future.delayed(widget.interval, _connect);
  }

  void _connect() {
    _subscription?.cancel();
    widget
        .watch(
          current.media.id,
          options: [authn.request(authn.AuthzCache.meta(context))],
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
        .then((_) {
          ds.RefreshBoundary.of(context)?.reset();
          debugPrint('download watch stream closed cleanly, reconnecting');
          _reconnect();
        })
        .catchError((e) {
          debugPrint('download watch socket closed by server, reconnecting: $e');
          _reconnect();
        }, test: ds.ErrorTests.socketclosed)
        .catchError((e) {
          debugPrint('download watch websocket closed abnormally, reconnecting: $e');
          _reconnect();
        }, test: ds.ErrorTests.websocketclosed)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: _resetcause);
          });
        });
  }

  @override
  void initState() {
    super.initState();
    current = widget.current;
    _connect();
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
        help: ds.Hint.multiline([
          Text("An active download showing progress, peer count, and completion percentage."),
          ds.HelpLabelled(
            label: Text("pause"),
            description: Text("suspend the download"),
          ),
          ds.HelpLabelled(
            label: Text("check"),
            description: Text("mark as processed once completed"),
          ),
        ]),
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
  final Widget help;
  const DownloadRowDisplay({super.key, required this.current, this.trailing, this.help = ds.HelpScope.None});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    var percentage = math.min(
      current.bytes == 0 ? 0.0 : (current.downloaded.toDouble() / current.bytes.toDouble()),
      1.0,
    );

    return ds.Help(
      ds.ErrorBoundary(
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
                        Expanded(
                          child: Text(
                            current.peers.toString().padLeft(3),
                            style: const TextStyle(fontFamily: 'monospace'),
                          ),
                        ),
                        Expanded(
                          child: Text(
                            "${(percentage * 100).toStringAsFixed(2).padLeft(6)}%",
                            style: const TextStyle(fontFamily: 'monospace'),
                          ),
                        ),
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
                  Text(current.peers.toString().padLeft(3), style: const TextStyle(fontFamily: 'monospace')),
                  Text(
                    "${(percentage * 100).toStringAsFixed(2).padLeft(6)}%",
                    style: const TextStyle(fontFamily: 'monospace'),
                  ),
                  trailing?.call(context) ?? const SizedBox(),
                ],
              );
            },
          ),
        ),
      ),
      help,
    );
  }
}
