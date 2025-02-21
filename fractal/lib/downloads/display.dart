import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import 'downloading.list.dart';
import 'available.list.dart';

class Display extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final defaults = ds.theme(context);
    return Container(
      padding: defaults.padding,
      child: Column(
        children: [
          DownloadingListDisplay(),
          Expanded(child: AvailableListDisplay()),
        ],
      ),
    );
  }
}
