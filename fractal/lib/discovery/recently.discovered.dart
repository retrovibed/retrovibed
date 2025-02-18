import 'package:flutter/material.dart';
import 'card.display.dart';

class Discovered extends StatelessWidget {
  const Discovered({super.key});

  @override
  Widget build(BuildContext context) {
    return GridView.count(
      primary: false,
      // padding: const EdgeInsets.all(20),
      crossAxisSpacing: 10,
      mainAxisSpacing: 10,
      crossAxisCount: 3,
      children: const <Widget>[
        CardDisplay(display: 'Discovered 1'),
        CardDisplay(display: 'Discovered 2'),
        CardDisplay(display: 'Discovered 3'),
      ],
    );
  }
}
