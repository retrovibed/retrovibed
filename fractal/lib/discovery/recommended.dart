import 'package:flutter/material.dart';
import 'card.display.dart';

class Recommended extends StatelessWidget {
  const Recommended({super.key});

  @override
  Widget build(BuildContext context) {
    return GridView.count(
      primary: false,
      crossAxisSpacing: 10,
      mainAxisSpacing: 10,
      crossAxisCount: 3,
      children: const <Widget>[
        CardDisplay(display: 'Recommended 1'),
        CardDisplay(display: 'Recommended 2'),
        CardDisplay(display: 'Recommended 3'),
      ],
    );
  }
}
