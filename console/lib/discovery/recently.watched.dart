import 'package:flutter/material.dart';
import 'card.display.dart';

class Watched extends StatelessWidget {
  const Watched({super.key});

  @override
  Widget build(BuildContext context) {
    // return Flex(
    //   direction: Axis.horizontal,
    //   children: const <Widget>[
    //     CardDisplay(display: 'Recent 1'),
    //     CardDisplay(display: 'Recent 2'),
    //     CardDisplay(display: 'Recent 3'),
    //   ],
    // );
    return Expanded(
      child: GridView.count(
        primary: false,
        crossAxisSpacing: 10,
        mainAxisSpacing: 10,
        crossAxisCount: 3,
        children: const <Widget>[
          CardDisplay(display: 'Recent 1'),
          CardDisplay(display: 'Recent 2'),
          CardDisplay(display: 'Recent 3'),
        ],
      ),
    );
  }
}
