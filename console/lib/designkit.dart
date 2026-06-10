import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/modals.dart' as _modals;
import 'design.kit/flutterx.dart';

export 'package:fixnum/fixnum.dart';
export 'design.kit/screens.dart';
export 'design.kit/accordian.dart';
export 'design.kit/buttons.dart';
export 'design.kit/buttons.loading.icon.dart';
export 'design.kit/buttons.loading.widget.dart';
export 'design.kit/carousel.dart';
export 'design.kit/card.dart';
export 'design.kit/confirmation.dart';
export 'design.kit/container.dart';
export 'design.kit/debug.dart';
export 'design.kit/bytesx.dart';
export 'design.kit/empty.dart';
export 'design.kit/errors.dart';
export 'design.kit/file.drop.well.dart';
export 'design.kit/future.dart';
export 'design.kit/guage.dart';
export 'design.kit/grid.dart';
export 'design.kit/heading.dart';
export 'design.kit/help.dart';
export 'design.kit/hyperlink.dart';
export 'design.kit/image.dart';
export 'design.kit/key.press.aware.dart';
export 'design.kit/long.hold.dart';
export 'design.kit/modals.dart';
export 'design.kit/periodic.dart';
export 'design.kit/refresh.dart';
export 'design.kit/repeat.dart';
export 'design.kit/rating.dart';
export 'design.kit/shake.dart';
export 'design.kit/shortcuts.dart';
export 'design.kit/search.dropdown.dart';
export 'design.kit/compacting.menu.dart';
export 'design.kit/search.tray.dart';
export 'design.kit/tables.dart';
export 'design.kit/theme.defaults.dart';
export 'design.kit/typography.dart';
export 'design.kit/flutterx.dart' show postframe;

abstract class modals {
  static _modals.NodeState? of(BuildContext context) {
    return _modals.of(context);
  }

  static void push(BuildContext context, Widget? child) {
    of(context)?.push(child);
  }

  static Future<T> asyncfn<T>(
    BuildContext context,
    Widget Function(Completer<T> completion) builder,
  ) => _modals.asyncfn(context, builder);
}

abstract class textediting {
  static void refocus(TextEditingController? controller) {
    if (controller == null) return;
    postframe(() {
      controller.selection = TextSelection.fromPosition(
        TextPosition(offset: controller.text.length),
      );
    });
  }
}

Widget build(Widget Function(BuildContext) b) {
  return Builder(builder: b);
}

Widget layout(Widget Function(BuildContext, BoxConstraints) b) {
  return LayoutBuilder(builder: b);
}

VoidCallback once(VoidCallback action) {
  bool hasRun = false;

  return () {
    if (hasRun) return;
    action();
    hasRun = true;
  };
}

// The target non-nullable asynchronous function type
typedef AsyncVoidCallback = Future<void> Function();

AsyncVoidCallback toasync(VoidCallback fn) {
  return () async {
    fn();
  };
}

AsyncVoidCallback decorated(AsyncVoidCallback fn, Future<void> Function(Future<void>) decorate) {
  return () => decorate(fn());
}
