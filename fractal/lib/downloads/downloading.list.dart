import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import 'package:fractal/media.dart' as media;

class DownloadingListDisplay extends StatelessWidget {
  final media.FnDownloadSearch search;
  const DownloadingListDisplay({
    super.key,
    this.search = media.discovered.downloading,
  });

  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
      initialData: <Widget>[],
      future: search(media.discoveredsearch.request(limit: 3)).then(
        (v) =>
            v.items.map((v) => media.DownloadRowDisplay(current: v)).toList(),
      ),
      builder: (BuildContext ctx, AsyncSnapshot<List<Widget>> snapshot) {
        if (snapshot.hasError) {
          return ds.Error.unknown(snapshot.error!);
        }

        return ds.Loading(
          loading: snapshot.connectionState != ConnectionState.done,
          child: ListView(shrinkWrap: true, children: snapshot.data!),
        );
      },
    );
  }
}
