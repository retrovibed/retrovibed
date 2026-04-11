import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;

class CardDisplay extends StatelessWidget {
  final String display;
  const CardDisplay({super.key, this.display = ''});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return Container(
      width: 128,
      height: 128,
      margin: defaults.margin,
      decoration: BoxDecoration(border: defaults.border),
      child: Flex(
        direction: Axis.vertical,
        children: [
          const Expanded(child: const Icon(Icons.movie, size: 128)),
          const Spacer(),
          Text(display),
        ],
      ),
    );
  }
}
