import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as media;

class DownloadingListDisplay extends StatefulWidget {
  final media.FnDownloadSearch search;
  final media.FnDownloadWatch watch;
  final ValueNotifier<int>? events;
  const DownloadingListDisplay({
    super.key,
    this.search = media.discovered.downloading,
    this.watch = media.discovered.watch,
    this.events,
  });

  @override
  State<StatefulWidget> createState() => _DownloadingListState();
}

class _DownloadingListState extends State<DownloadingListDisplay> {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  Timer? period;
  media.DownloadSearchResponse _res = media.discoveredsearch.response(
    next: media.discoveredsearch.request(limit: 3),
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

  void refresh() {
    widget
        .search(
          _res.next,
          options: [authn.request(authn.AuthzCache.meta(context))],
        )
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
    ds.postframe(() => refresh());
    period = Timer.periodic(
      const Duration(seconds: 20),
      (p) => refresh(),
    );
    widget.events?.addListener(() {
      refresh();
    });
  }

  @override
  void dispose() {
    super.dispose();
    period?.cancel();
  }

  @override
  Widget build(BuildContext context) {
    return ds.RefreshBoundary(
      onReset: () {
        widget.events ?? refresh();
        widget.events?.value += 1;
      },
      ds.Table(
        loading: _loading,
        cause: _cause,
        children: _res.items,
        ds.Table.inline<media.Download>(
          (v) => ds.ErrorBoundary(
            media.RefreshingDownload(
              key: ValueKey(v.media.id),
              current: v,
              watch: widget.watch,
            ),
          ),
        ),
      ),
    );
  }
}
