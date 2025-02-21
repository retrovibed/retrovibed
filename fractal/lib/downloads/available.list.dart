import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import 'package:fractal/media.dart' as media;

class AvailableListDisplay extends StatelessWidget {
  final media.FnMediaSearch search;
  const AvailableListDisplay({
    super.key,
    this.search = media.discovered.available,
  });

  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
      initialData: <Widget>[],
      future: search(media.mediasearch.request(limit: 32))
          .then(
            (v) =>
                v.items
                    .map(
                      (v) => media.RowDisplay(
                        media: v,
                        onTap:
                            () => media.discovered
                                .download(v.id)
                                .then((v) {
                                  ds.RefreshBoundary.of(context)?.reset();
                                })
                                .catchError((cause) {
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    SnackBar(
                                      content: Text(
                                        "Failed to download: $cause",
                                      ),
                                    ),
                                  );
                                  return null;
                                }),
                      ),
                    )
                    .toList(),
          )
          .catchError(
            ds.Error.boundary(
              context,
              List<media.RowDisplay>.empty(),
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
