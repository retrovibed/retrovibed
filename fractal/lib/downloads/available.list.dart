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
      future: search(media.mediasearch.request(limit: 32)).then(
        (v) =>
            v.items
                .map(
                  (v) => media.RowDisplay(
                    media: v,
                    onTap: () => print("download not yet implemented"),
                  ),
                )
                .toList(),
      ),
      builder: (BuildContext ctx, AsyncSnapshot<List<Widget>> snapshot) {
        if (snapshot.hasError) {
          return ds.Error.unknown(snapshot.error!);
        }

        return ds.Loading(
          loading: snapshot.connectionState != ConnectionState.done,
          child: ListView(children: snapshot.data!),
        );
      },
    );
  }
}
