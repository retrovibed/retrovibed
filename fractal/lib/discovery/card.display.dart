import 'package:flutter/material.dart';

class CardDisplay extends StatelessWidget {
  final String display;
  const CardDisplay({super.key, this.display = ''});

  @override
  Widget build(BuildContext context) {
    // return Flex(
    //   direction: Axis.vertical,
    //   children: [
    //     const Expanded(child: const Icon(Icons.movie, size: 128)),
    //     const Spacer(),
    //     Text(display),
    //   ],
    // );
    return Container(
      width: 128,
      height: 128,
      decoration: const BoxDecoration(
        border: Border(
          top: BorderSide(color: Color(0xFF000000)),
          left: BorderSide(color: Color(0xFF000000)),
          right: BorderSide(color: Color(0xFF000000)),
          bottom: BorderSide(color: Color(0xFF000000)),
        ),
      ),
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
