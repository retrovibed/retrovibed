import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as media;
import 'download.watch.dart';

class DownloadingListDisplay extends StatefulWidget {
  final media.FnDownloadSearch search;
  final media.FnDownloadWatch watch;
  final StreamController<media.Download> events;
  const DownloadingListDisplay(
    this.events, {
    super.key,
    this.search = media.discovered.downloading,
    this.watch = media.discovered.watch,
  });

  @override
  State<StatefulWidget> createState() => _DownloadingListState();
}

class _DownloadingListState extends State<DownloadingListDisplay> with ds.LoadingState {
  Timer? period;
  StreamSubscription<void>? subscription;
  media.DownloadSearchResponse _res = media.discoveredsearch.response(
    next: media.discoveredsearch.request(limit: 3),
  );

  void refresh() {
    widget
        .search(
          _res.next,
          options: [authn.request(authn.AuthzCache.meta(context))],
        )
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
    ds.postframe(() => refresh());
    period = Timer.periodic(
      const Duration(seconds: 20),
      (p) => refresh(),
    );
    subscription = widget.events.stream.listen((v) {
      // only an empty download signals a refresh.
      if (v.media.id == "") return refresh();

      setState(() {
        _res = media.DownloadSearchResponse(
          items: ds.fnOnChangeOrInsert(_res.items, v, (d) => d.media.id == v.media.id),
          next: _res.next,
        );
      });
    });
  }

  @override
  void dispose() {
    super.dispose();
    period?.cancel();
    subscription?.cancel();
  }

  @override
  Widget build(BuildContext context) {
    return ds.RefreshBoundary(
      onReset: () {
        widget.events.add(media.Download());
      },
      ds.Table(
        loading: loading,
        cause: cause,
        children: _res.items,
        collapsable: true,
        ds.Table.inline<media.Download>(
          (v) => ds.ErrorBoundary(
            RefreshingDownload(
              key: ValueKey(v.media.id),
              current: v,
              watch: widget.watch,
              updates: widget.events,
              onCompleted: (d) {
                refresh();
              },
            ),
          ),
        ),
      ),
    );
  }
}
