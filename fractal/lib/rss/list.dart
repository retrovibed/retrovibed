import 'package:flutter/material.dart';
import 'package:fractal/rss.dart';

class Item extends StatelessWidget {
  final Feed current;
  const Item({super.key, required Feed this.current});

  @override
  Widget build(BuildContext context) {
    return Container(child: Edit(feed: this.current));
  }
}

class ListFeeds extends StatelessWidget {
  final List<Feed> current;
  const ListFeeds({super.key, required List<Feed> this.current});

  @override
  Widget build(BuildContext context) {
    if (this.current.isEmpty) {
      return Container();
    }

    return Container(
      child: Column(
        spacing: 5.0,
        children: this.current.map((f) => Item(current: f)).toList(),
      ),
    );
  }
}
