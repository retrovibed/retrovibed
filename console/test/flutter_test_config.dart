import 'dart:async';
import 'dart:io';

import 'package:flutter/foundation.dart';

Future<void> testExecutable(FutureOr<void> Function() testMain) {
  // Only silence output when stdout isn't a terminal (e.g. piped to a file/
  // log collector) or we're running under CI (which often still reports a
  // terminal). When run directly in an interactive shell, leave
  // print/debugPrint alone so local debugging still works.
  final isCI = Platform.environment['CI'] != null;
  if (!stdout.hasTerminal || isCI) {
    // Flutter's own logging function (used by framework code and anyone
    // calling debugPrint directly) routes through this hook, so replacing
    // it with a no-op silences it for every test in this directory tree.
    debugPrint = (String? message, {int? wrapWidth}) {};

    // Plain print() doesn't go through debugPrint, so it has to be caught
    // separately by running tests inside a zone whose print handler is a
    // no-op instead of the default (which forwards to stdout).
    return runZoned(
      () async => testMain(),
      zoneSpecification: ZoneSpecification(
        print: (self, parent, zone, line) {},
      ),
    );
  }

  return Future.sync(testMain);
}
