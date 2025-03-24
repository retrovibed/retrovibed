import 'package:flutter/material.dart';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:console/designkit.dart' as ds;

class SearchTray extends StatelessWidget {
  final Widget trailing;
  final void Function(String i) onSubmitted;
  final void Function(fixnum.Int64 i) next;
  final fixnum.Int64 current;
  final bool empty;
  final bool autofocus;
  final TextEditingController? controller;
  final FocusNode? focus;

  SearchTray({
    super.key,
    required this.onSubmitted,
    required this.next,
    required this.current,
    required this.empty,
    this.trailing = const SizedBox(),
    this.autofocus = false,
    this.controller,
    this.focus,
  });

  @override
  Widget build(BuildContext context) {
    final theming = ds.Defaults.of(context);
    return Container(
      padding: theming.padding,
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Expanded(
            child: TextField(
              controller: controller,
              decoration: InputDecoration(hintText: "search the library"),
              autofocus: autofocus,
              focusNode: focus,
              onSubmitted: onSubmitted,
            ),
          ),
          IconButton(
            onPressed: current == 0 ? null : () => next(current + 1),
            icon: Icon(Icons.arrow_left),
          ),
          IconButton(
            onPressed: empty ? null : () => next(current + 1),
            icon: Icon(Icons.arrow_right),
          ),
          trailing,
        ],
      ),
    );
  }
}
