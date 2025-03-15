import 'package:flutter/material.dart';
import 'package:fractal/rss.dart';
import './api.dart' as api;

void _Noop(Feed up) {}

class Item extends StatelessWidget {
  final Feed current;
  final void Function(Feed upd) onChange;
  const Item({super.key, required Feed this.current, this.onChange = _Noop});

  @override
  Widget build(BuildContext context) {
    return Container(
      child: Edit(
        feed: this.current,
        onChange: (u) {
          api.create(FeedCreateRequest(feed: u)).then((resp) {
            onChange(resp.feed);
          });
        },
      ),
    );
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

    // return Container();
    return Container(
      child: Column(
        spacing: 5.0,
        children: this.current.map((f) => Item(current: f)).toList(),
      ),
    );
  }
}
