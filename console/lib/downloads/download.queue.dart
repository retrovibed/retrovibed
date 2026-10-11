import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media/api.dart' as api;
import 'download.watch.dart';

class DownloadQueue extends StatefulWidget {
  final ValueNotifier<Future<List<api.Download>>> queue;
  final Duration interval;
  final api.FnDownloadSearch apidownloading;
  final api.FnDownloadWatch apiwatch;
  final Duration minCompletedDisplay;
  const DownloadQueue(
    this.queue, {
    super.key,
    this.interval = const Duration(milliseconds: 5000),
    this.apidownloading = api.discovered.downloading,
    this.apiwatch = api.discovered.watch,
    this.minCompletedDisplay = const Duration(seconds: 15),
  });

  @override
  State<DownloadQueue> createState() => _DownloadQueue();
}

class _DownloadQueue extends State<DownloadQueue> with ds.LoadingState {
  List<api.Download> _pending = [];
  final StreamController<api.Download> _updates = StreamController<api.Download>.broadcast();
  StreamSubscription<api.Download> _subscription = const Stream<api.Download>.empty().listen(null);

  @override
  void dispose() {
    widget.queue.removeListener(_load);
    _subscription.cancel();
    _updates.close();
    super.dispose();
  }

  void _onItemCompleted(api.Download startedAs) {
    final delay = api.download.completed(startedAs) ? widget.minCompletedDisplay : Duration.zero;
    Future.delayed(delay, () {
      setState(() {
        _pending = ds.fnOnChange(_pending, null, (d) => d.media.id == startedAs.media.id);
      });
      // queue drained, check for any remaining active downloads.
      if (_pending.isEmpty) _load();
    });
  }

  void _load() {
    setState(() {
      loading = true;
    });
    widget.queue.value
        .then(
          (_) => widget.apidownloading(
            api.discoveredsearch.request(limit: 128),
            options: [authn.request(authn.AuthzCache.meta(context))],
          ),
        )
        .then((v) => v.items)
        .then((list) {
          setState(() {
            _pending = list;
            loading = false;
          });
        })
        .catchError((cause) {
          setState(() {
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
            loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();
    widget.queue.addListener(_load);
    _subscription = _updates.stream.listen((v) {
      // completion is handled by onCompleted, replacing the item here would change the completed display hold.
      if (api.download.completed(v)) return;
      if (api.download.paused(v)) return _load();
      setState(() {
        _pending = ds.fnOnChange(_pending, v, (d) => d.media.id == v.media.id);
      });
    });
    ds.postframe(() => _load());
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ds.Loading(
      loading: loading,
      cause: cause,
      ds.ErrorBoundary(
        _pending.isEmpty
            ? const SizedBox.shrink()
            : Builder(
                builder: (context) {
                  final item = _pending.first;
                  return ds.Container(
                    padding: defaults.padding,
                    RefreshingDownload(
                      key: ValueKey("${item.media.id}.${item.media.updatedAt}"),
                      current: item,
                      interval: widget.interval,
                      watch: widget.apiwatch,
                      updates: _updates,
                      onCompleted: (_) => _onItemCompleted(item),
                    ),
                  );
                },
              ),
      ),
    );
  }
}
