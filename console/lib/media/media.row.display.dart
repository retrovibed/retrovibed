import 'package:flutter/material.dart';
import 'package:console/designkit.dart' as ds;
import './media.pb.dart';

class RowDisplay extends StatelessWidget {
  final Media media;
  final List<Widget> leading;
  final List<Widget> trailing;
  final void Function()? onTap;
  const RowDisplay({
    super.key,
    required this.media,
    this.leading = const [],
    this.trailing = const [],
    this.onTap = ds.defaulttap,
  });

  @override
  Widget build(BuildContext context) {
    final themex = ds.Defaults.of(context);
    List<Widget> children = List.from(leading);
    children += [
      Expanded(child: Text(media.description, overflow: TextOverflow.ellipsis)),
    ];
    children += trailing;

    return Container(
      padding: themex.padding,
      child: InkWell(
        onTap: onTap,
        child: Row(spacing: themex.spacing!, children: children),
      ),
    );
  }
}
