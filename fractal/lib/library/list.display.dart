import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import 'package:fractal/media.dart' as media;
import './search.row.dart';
import 'file.drop.well.dart';

class AvailableListDisplay extends StatefulWidget {
  final media.FnMediaSearch search;
  const AvailableListDisplay({super.key, this.search = media.mediasearch.get});

  @override
  State<StatefulWidget> createState() => _AvailableListDisplay();
}

class _AvailableListDisplay extends State<AvailableListDisplay> {
  bool _loading = true;
  ds.Error? _cause = null;
  media.MediaSearchResponse _res = media.mediasearch.response(
    next: media.mediasearch.request(limit: 32),
  );

  void refresh() {
    widget
        .search(_res.next)
        .then((v) {
          setState(() {
            _res = v;
            _loading = false;
          });
        })
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e);
            _loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();
    refresh();
  }

  @override
  Widget build(BuildContext context) {
    return ds.Table(
      loading: _loading,
      cause: _cause,
      leading: SearchTray(
        onSubmitted: (v) {
          setState(() {
            _res.next.query = v;
            _res.next.offset = fixnum.Int64(0);
          });
          refresh();
        },
        next: (i) {
          setState(() {
            _res.next.offset = i;
          });
          refresh();
        },
        current: _res.next.offset,
        empty: fixnum.Int64(_res.items.length) < _res.next.limit,
      ),
      children: _res.items,
      (v) => media.RowDisplay(media: v),
      empty: FileDropWell((v) {
        return Future.value(null);
      }),
    );
  }
}
