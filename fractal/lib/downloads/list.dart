import 'package:flutter/material.dart';
import 'horizontal.display.dart';

class List extends StatelessWidget {
  const List({super.key});

  @override
  Widget build(BuildContext context) {
    return ListView(
      children: const <Widget>[
        HorizontalDisplay(
          display: 'Download 1',
        ),
        HorizontalDisplay(
          display: 'Download 2',
        ),
        HorizontalDisplay(
          display: 'Download 3',
        ),
      ],
    );
  }
}
