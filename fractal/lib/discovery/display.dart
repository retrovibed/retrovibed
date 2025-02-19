import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import 'card.display.dart';

Future<List<Widget>> data() {
  var f1 = Future.delayed(
    Duration(seconds: 1),
    () => <Widget>[
      CardDisplay(display: 'Recent 1'),
      CardDisplay(display: 'Recent 2'),
      CardDisplay(display: 'Recent 3'),
    ],
  );

  var f2 = Future.delayed(
    Duration(seconds: 1),
    () => <Widget>[
      CardDisplay(display: 'Discovered 1'),
      CardDisplay(display: 'Discovered 2'),
      CardDisplay(display: 'Discovered 3'),
    ],
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
        if (snapshot.hasError) {
          return SizedBox.expand(child: Text("failed"));
        }

        return SizedBox.expand(
          child: ds.Loading(
            loading: snapshot.connectionState != ConnectionState.done,
            child: GridView.count(
              primary: false,
              padding: const EdgeInsets.all(20),
              crossAxisSpacing: 0,
              mainAxisSpacing: 0,
              crossAxisCount: 3,
              children: snapshot.data!,
            ),
          ),
        );
      },
    );
  }
}
