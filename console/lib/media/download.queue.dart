import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as api;
import 'download.watch.dart';
import 'package:retrovibed/design.kit/stateful.dart';

class DownloadQueue extends StatefulWidget {
  final Future<List<api.Download>> queue;
  final Duration interval;
  final api.FnDownloadWatch watch;
  final Duration minCompletedDisplay;
  final void Function()? onQueueComplete;
  const DownloadQueue(
    this.queue, {
    super.key,
    this.interval = const Duration(milliseconds: 5000),
    this.watch = api.discovered.watch,
    this.minCompletedDisplay = const Duration(seconds: 15),
    this.onQueueComplete,
  });

  @override
  State<DownloadQueue> createState() => _DownloadQueue();
}

class _DownloadQueue extends State<DownloadQueue> with LoadingState {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  List<api.Download> _pending = [];
  int _index = 0;

  void _resetcause() {
    setState(() => _cause = ds.Error.zero);
  }

  void _show(int index) {
    setState(() => _index = index);
    if (index >= _pending.length) {
      widget.onQueueComplete?.call();
    }
  }

  void _advance(int index) {
    if (!mounted || index != _index) return;
    _show(index + 1);
  }

  void _onItemCompleted(int index, api.Download startedAs) {
    final delay = api.download.completed(startedAs) ? widget.minCompletedDisplay : Duration.zero;
    Future.delayed(delay, () => _advance(index));
  }

  @override
  void initState() {
    super.initState();
    widget.queue
        .then((list) {
          if (!mounted) return;
          setState(() {
            _pending = list;
            _loading = false;
          });
          _show(0);
        })
        .catchError((cause) {
          if (!mounted) return;
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: _resetcause);
            _loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ds.Loading(
      loading: _loading,
      cause: _cause,
      ds.ErrorBoundary(
        _index >= _pending.length
            ? const SizedBox.shrink()
            : Builder(
                builder: (context) {
                  final item = _pending[_index];
                  return ds.Container(
                    padding: defaults.padding,
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      spacing: defaults.spacing,
                      children: [
                        Text(
                          "${_index + 1} of ${_pending.length}",
                          style: const TextStyle(fontFamily: 'monospace'),
                        ),
                        RefreshingDownload(
                          key: ValueKey(item.media.id),
                          current: item,
                          interval: widget.interval,
                          watch: widget.watch,
                          onCompleted: (_) => _onItemCompleted(_index, item),
                        ),
                      ],
                    ),
                  );
                },
              ),
      ),
    );
  }
}
