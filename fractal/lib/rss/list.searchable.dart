import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import 'package:fractal/rss.dart';
import 'package:fixnum/fixnum.dart';
import 'package:fractal/rss/list.dart';
import './api.dart' as api;

class ListSearchable extends StatefulWidget {
  final api.FnSearch search;

  ListSearchable({super.key, this.search = api.search});

  @override
  State<ListSearchable> createState() => SearchableView();
}

class SearchableView extends State<ListSearchable> {
  FeedSearchResponse current = FeedSearchResponse(
    next: FeedSearchRequest(query: '', offset: Int64(0), limit: Int64(10)),
    items: [],
  );
  Future<FeedSearchResponse> pending = Future.delayed(
    Duration(hours: 999999),
    () => FeedSearchResponse(
      next: FeedSearchRequest(query: '', offset: Int64(0), limit: Int64(10)),
      items: [],
    ),
  );

  Future<FeedSearchResponse> refresh() {
    return widget.search(current.next).then((r) {
      setState(() {
        current = r;
      });
      return r;
    });
  }

  @override
  void initState() {
    super.initState();
    pending = refresh();
  }

  @override
  Widget build(BuildContext context) {
    final zerolist = FeedSearchResponse(
      items: List.generate(current.next.limit.toInt(), (idx) => Feed.create()),
    );

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Expanded(
              child: TextField(
                decoration: InputDecoration(hintText: "search feeds"),
                onChanged:
                    (v) => setState(() {
                      current.next.query = v;
                    }),
                onSubmitted:
                    (v) => setState(() {
                      pending = refresh();
                    }),
              ),
            ),
            IconButton(onPressed: () {}, icon: Icon(Icons.arrow_left)),
            IconButton(onPressed: () {}, icon: Icon(Icons.arrow_right)),
          ],
        ),
        FutureBuilder(
          initialData: zerolist,
          future: pending,
          builder: (
            BuildContext ctx,
            AsyncSnapshot<FeedSearchResponse> snapshot,
          ) {
            if (snapshot.hasError) {
              print(snapshot.error);
              return SizedBox.expand(child: Text("failed"));
            }

            return ds.Loading(
              loading: snapshot.connectionState != ConnectionState.done,
              child: ListFeeds(current: current.items),
            );
          },
        ),
      ],
    );
  }
}
