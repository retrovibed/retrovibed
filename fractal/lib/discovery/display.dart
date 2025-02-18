import 'package:flutter/material.dart';
import 'card.display.dart';

class Display extends StatelessWidget {
  const Display({super.key});

  @override
  Widget build(BuildContext context) {
    return SizedBox.expand(
      child: GridView.count(
        primary: false,
        padding: const EdgeInsets.all(20),
        crossAxisSpacing: 0,
        mainAxisSpacing: 0,
        crossAxisCount: 3,
        children: const <Widget>[
          CardDisplay(display: 'Recent 1'),
          CardDisplay(display: 'Recent 2'),
          CardDisplay(display: 'Recent 3'),
          CardDisplay(display: 'Discovered 1'),
          CardDisplay(display: 'Discovered 2'),
          CardDisplay(display: 'Discovered 3'),
          CardDisplay(display: 'Recommended 1'),
          CardDisplay(display: 'Recommended 2'),
          CardDisplay(display: 'Recommended 3'),
        ],
      ),
    );
  }
}
