import 'package:flutter/material.dart';
import 'package:fractal/discovery/recently.discovered.dart';
import 'package:fractal/discovery/recently.watched.dart';
import 'package:fractal/discovery/recommended.dart';

class Display extends StatelessWidget {
  const Display({super.key});

  @override
  Widget build(BuildContext context) {
    return const Column(
      children: [
        RecentlyWatched(),
        Discovered(),
        Recommended(),
      ],
    );
  }
}
