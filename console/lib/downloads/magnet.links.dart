import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;

class MagnetDownloads extends StatelessWidget {
  final FocusNode _focus = FocusNode();
  final TextEditingController controller;
  final Future<void> Function(List<String>) onSubmitted;
  final StreamController<String> stream = StreamController<String>();

  MagnetDownloads({super.key, controller, required this.onSubmitted})
    : this.controller = controller ?? TextEditingController();

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return Container(
      padding: defaults.padding,
      color: defaults.opaque,
      child: Column(
        spacing: defaults.spacing,
        mainAxisSize: MainAxisSize.min,
        children: [
          TextField(
            focusNode: _focus,
            controller: controller,
            decoration: const InputDecoration(
              labelText: 'Enter an magnet link(s)',
              border: OutlineInputBorder(),
            ),
            onSubmitted: (v) {
              stream.add(v);
              controller.clear();
              ds.textediting.refocus(controller);
              _focus.requestFocus();
            },
          ),
          ds.ErrorBoundary(
            ds.build((context) {
              return forms.ItemListManager<String>(
                stream: stream.stream,
                onSubmitted: (l) {
                  onSubmitted([
                    if (controller.text.isNotEmpty) controller.text,
                    ...l,
                  ]).catchError(
                    ds.Error.boundary(context, null, ds.Error.unknown),
                  );
                },
                builder: (item) {
                  return SelectionArea(
                    child: Text(
                      item,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  );
                },
              );
            }),
          ),
        ],
      ),
    );
  }
}
