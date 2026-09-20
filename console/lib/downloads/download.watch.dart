import 'dart:async';
import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media/api.dart' as api;
import 'download.row.controls.dart';
import 'download.row.display.dart';

class RefreshingDownload extends StatefulWidget {
  final api.Download current;
  final Duration interval;
  final api.FnDownloadWatch watch;
  final Function(api.Download) onCompleted;

  /// receives the download each time it changes.
  final StreamSink<api.Download> updates;

  /// quiet period required before the latest update is pushed to updates.
  final Duration debounce;
  const RefreshingDownload({
    super.key,
    required this.current,
    this.interval = const Duration(milliseconds: 5000),
    this.watch = api.discovered.watch,
    this.onCompleted = ds.fnNoop,
    required this.updates,
    this.debounce = const Duration(milliseconds: 250),
  });

  @override
  State<RefreshingDownload> createState() => _DownloadingState();
}

class _DownloadingState extends State<RefreshingDownload> with ds.LoadingState {
  api.Download current = api.Download();
  StreamSubscription<api.Download> _subscription = const Stream<api.Download>.empty().listen(null);
  Timer _debounce = Timer(Duration.zero, () {});
  bool _notifiedCompleted = false;

  /// transfer rate in bytes per second, derived from successive updates.
  int rate = 0;
  DateTime _sampled = DateTime.now();
  int _sampledBytes = 0;

  /// minimum period between rate updates.
  static const Duration _rateInterval = Duration(seconds: 1);

  void _sample(api.Download v) {
    final now = DateTime.now();
    final elapsed = now.difference(_sampled);
    if (elapsed < _rateInterval) return;

    rate = math.max(0, ((v.downloaded.toInt() - _sampledBytes) * Duration.microsecondsPerSecond / elapsed.inMicroseconds).round());
    _sampled = now;
    _sampledBytes = v.downloaded.toInt();
  }

  void _maybeNotifyCompleted() {
    if (_notifiedCompleted || !api.download.completed(current)) return;
    print("notified completed ${_notifiedCompleted} ${api.download.completed(current)} ${current}");
    _notifiedCompleted = true;
    widget.onCompleted(current);
  }

  void _reconnect() {
    if (!mounted) return;
    Future.delayed(widget.interval, _connect);
  }

  void _connect() {
    _subscription.cancel();
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
                _sample(v);
                current = v;
              });
              _debounce.cancel();
              _debounce = Timer(widget.debounce, () => widget.updates.add(v));
              ds.postframe(_maybeNotifyCompleted);
            },
            cancelOnError: true,
            onError: c.completeError,
            onDone: c.complete,
          );

          return c.future;
        })
        .then((x) {
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
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
          });
        });
  }

  @override
  void initState() {
    super.initState();
    current = widget.current;
    _sampledBytes = current.downloaded.toInt();
    ds.postframe(_maybeNotifyCompleted);
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    _connect();
  }

  @override
  void dispose() {
    _debounce.cancel();
    _subscription.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ds.ErrorScreen(
      cause: cause,
      DownloadRowDisplay(
        current: current,
        rate: rate,
        help: ds.Hint.multiline([
          Text("An active download showing progress, peer count, transfer rate, and completion percentage."),
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
