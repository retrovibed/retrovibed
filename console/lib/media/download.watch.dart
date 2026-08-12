import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'api.dart' as api;
import 'download.row.controls.dart';
import 'download.row.display.dart';
import 'package:retrovibed/design.kit/stateful.dart';

class RefreshingDownload extends StatefulWidget {
  static _noopOnCompleted(api.Download v) {}
  final api.Download current;
  final Duration interval;
  final api.FnDownloadWatch watch;
  final Function(api.Download) onCompleted;
  const RefreshingDownload({
    super.key,
    required this.current,
    this.interval = const Duration(milliseconds: 5000),
    this.watch = api.discovered.watch,
    this.onCompleted = RefreshingDownload._noopOnCompleted,
  });

  @override
  State<RefreshingDownload> createState() => _DownloadingState();
}

class _DownloadingState extends State<RefreshingDownload> with LoadingState {
  Widget _cause = ds.Error.zero;
  api.Download current = api.Download();
  StreamSubscription<api.Download>? _subscription;
  bool _notifiedCompleted = false;

  void _maybeNotifyCompleted() {
    if (_notifiedCompleted || !api.download.completed(current)) return;
    _notifiedCompleted = true;
    widget.onCompleted(current);
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
              ds.postframe(_maybeNotifyCompleted);
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
    ds.postframe(_maybeNotifyCompleted);
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (_subscription == null) {
      _connect();
    }
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
        trailing: (ctx) => DownloadRowControls(
          current: current,
          onChange: (d) {
            ds.RefreshBoundary.of(ctx)?.reset();
          },
        ),
      ),
    );
  }
}
