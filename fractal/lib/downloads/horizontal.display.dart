import 'package:flutter/material.dart';

class HorizontalDisplay extends StatelessWidget {
  final String display;
  const HorizontalDisplay({super.key, this.display = ''});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        const Icon(Icons.map),
        Text(display),
        const Spacer(),
        const Expanded(
          child: LinearProgressIndicator(
            value: 0.2,
            semanticsLabel: 'Linear progress indicator',
          ),
        ),
        const Spacer(),
        const Icon(Icons.delete)
      ],
    );
  }
}
