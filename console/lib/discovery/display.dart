import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'card.display.dart';
import 'package:retrovibed/media.dart' as media;

Future<List<Widget>> data() {
  // var f1 = media.recent().then(
  //   (v) => v.items.map((v) => CardDisplay(display: v.description)).toList(),
  // );
  var f1 = Future.delayed(
    Duration(seconds: 1),
    () => <Widget>[
      CardDisplay(display: 'Recent 1'),
      CardDisplay(display: 'Recent 2'),
      CardDisplay(display: 'Recent 3'),
    ],
  );

  var f2 = media.discovered
      .available(media.discoveredsearch.request())
      .then(
        (v) => v.items.map((v) => CardDisplay(display: v.media.description)).toList(),
      );

  var f3 = Future.delayed(
    Duration(seconds: 1),
    () => <Widget>[
      CardDisplay(display: 'Recommended 1'),
      CardDisplay(display: 'Recommended 2'),
      CardDisplay(display: 'Recommended 3'),
    ],
  );

  return Future.wait([f1, f2, f3]).then(
    (results) => results.reduce((sum, current) {
      sum.addAll(current);
      return sum;
    }),
  );
}

class Display extends StatelessWidget {
  const Display({super.key});

  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
      initialData: <Widget>[],
      future: data(),
      builder: (BuildContext ctx, AsyncSnapshot<List<Widget>> snapshot) {
        final defaults = ds.Defaults.of(context);

        if (snapshot.hasError) {
          print(snapshot.error);
          return ds.Error.text(snapshot.error.toString());
        }

        return ds.Loading(
          loading: snapshot.connectionState != ConnectionState.done,
          GridView.count(
            primary: false,
            padding: defaults.padding,
            crossAxisSpacing: 0,
            mainAxisSpacing: 0,
            crossAxisCount: 3,
            children: snapshot.data!,
          ),
        );
      },
    );
  }
}
