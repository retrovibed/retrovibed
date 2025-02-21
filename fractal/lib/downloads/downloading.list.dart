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
      future: search(media.discoveredsearch.request(limit: 3))
          .then(
            (v) =>
                v.items
                    .map(
                      (v) =>
                          media.DownloadRowDisplay(
                                current: v,
                                trailing:
                                    (ctx) => media.DownloadRowControls(
                                      current: v,
                                      onChange: (d) {
                                        ds.RefreshBoundary.of(ctx)?.reset();
                                      },
                                    ),
                              )
                              as Widget,
                    )
                    .toList(),
          )
          .catchError(
            ds.Error.boundary(
              context,
              List<media.DownloadRowDisplay>.empty(),
              ds.Error.offline,
            ),
            test: ds.ErrorTests.offline,
          )
          .catchError((e) => throw ds.Error.unknown(e)),
      builder: (BuildContext ctx, AsyncSnapshot<List<Widget>> snapshot) {
        return ds.Loading(
          loading: snapshot.connectionState != ConnectionState.done,
          cause: ds.Error.maybeErr(snapshot.error),
          child: ListView(shrinkWrap: true, children: snapshot.data ?? []),
        );
      },
    );
  }
}
