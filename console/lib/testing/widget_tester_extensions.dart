import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as meta;

bool _robotoLoaded = false;

Future<void> _loadRoboto() async {
  if (_robotoLoaded) return;
  // Resolve Roboto from the active Flutter SDK's material fonts cache.
  final flutterRoot =
      Platform.environment['FLUTTER_ROOT'] ?? File(Platform.resolvedExecutable).parent.parent.parent.path;
  final fontPath = '$flutterRoot/bin/cache/artifacts/material_fonts/Roboto-Regular.ttf';
  final file = File(fontPath);
  if (!file.existsSync()) return;
  final loader = FontLoader('Roboto')..addFont(Future.value(file.readAsBytesSync().buffer.asByteData()));
  await loader.load();
  _robotoLoaded = true;
}

extension WidgetTesterExtensions on WidgetTester {
  Future<void> pumpApp(
    Widget child, {
    ThemeData? theme,
    AlignmentGeometry alignment = Alignment.center,
    FlexFit fit = FlexFit.loose,
    Size physicalSize = const Size(800, 600),
    Future<meta.AuthzResponse> Function({String? host}) authzCurrent = authn.AuthzCache.fake,
  }) async {
    view.physicalSize = physicalSize;
    view.devicePixelRatio = 1.0;
    addTearDown(view.reset);
    await _loadRoboto();
    final base = theme ?? ThemeData();
    final defaultTextButtonStyle = TextButton.styleFrom(
      enabledMouseCursor: SystemMouseCursors.click,
    );
    final merged = base.copyWith(
      textButtonTheme: TextButtonThemeData(
        style: defaultTextButtonStyle.merge(base.textButtonTheme.style),
      ),
    );

    return pumpWidget(
      MaterialApp(
        theme: merged,
        home: Material(
          child: Align(
            alignment: alignment,
            child: Flex(
              direction: Axis.vertical,
              children: [
                Flexible(
                  fit: fit,
                  child: authn.AuthzCache(
                    LayoutBuilder(
                      builder: (context, constraints) {
                        final defaults = ds.Defaults.of(context);
                        final isCompact = constraints.maxWidth < defaults.compact;
                        return Theme(
                          data: Theme.of(context).copyWith(
                            extensions: [defaults.copyWith(isCompact: isCompact)],
                          ),
                          child: child,
                        );
                      },
                    ),
                    current: authzCurrent,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  /// Prints the full hit-test path at the center of [finder], mirroring
  /// exactly what the MouseTracker sees. Useful for debugging cursor and
  /// pointer-event issues in CI.
  void hitTestAudit(Finder finder) {
    final Offset point = getCenter(finder);
    final HitTestResult result = HitTestResult();
    binding.hitTestInView(result, point, this.view.viewId);

    print('--- Hit-Test Audit ---');
    for (final entry in result.path) {
      print('Hit: ${entry.target.runtimeType} -> ${entry.target}');
    }
  }

  /// Pumps frames until a widget matching [finder] satisfies [condition].
  Future<void> pumpUntil(
    Finder finder,
    bool Function(Widget widget) condition,
  ) async {
    int fetchedFrames = 0;
    while (!condition(widget(finder.first))) {
      await pump(const Duration(milliseconds: 16));
      fetchedFrames++;
      if (fetchedFrames > 120) {
        throw Exception(
          'Timed out waiting for condition on ${finder.describeMatch(Plurality.one)}',
        );
      }
    }
    await pumpAndSettle();
  }

  /// Returns the [MouseCursor] currently active for [device] as reported by
  /// [MouseTracker]. This reflects what Flutter has actually resolved for the
  /// pointer's current position, without needing a [Finder].
  ///
  /// [device] defaults to 1, which is the device id assigned by
  /// [WidgetTester.createGesture] for the first mouse pointer in a test.
  MouseCursor? currentCursor({int device = 1}) {
    return RendererBinding.instance.mouseTracker.debugDeviceActiveCursor(
      device,
    );
  }

  /// Returns the resolved [MouseCursor] at the center of [finder] by walking
  /// the hit-test path and returning the first non-deferred cursor found.
  MouseCursor resolvedCursorAt(Finder finder) {
    final Offset location = getCenter(finder);
    final HitTestResult result = HitTestResult();
    binding.hitTestInView(result, location, view.viewId);

    for (final entry in result.path) {
      final target = entry.target;
      if (target is RenderObject) {
        try {
          final dynamic potential = (target as dynamic).cursor;
          if (potential is MouseCursor && potential != MouseCursor.defer) {
            return potential;
          }
        } catch (_) {}
      }
    }

    return SystemMouseCursors.basic;
  }

  /// Pumps frames until the widget at [finder] is the actual hit target
  /// and is no longer blocked by active IgnorePointers or AbsorbPointers.
  Future<void> pumpUntilHitVisible(Finder finder) async {
    int frames = 0;
    bool blocked = true;

    while (blocked && frames < 120) {
      await pump(const Duration(milliseconds: 16));

      final HitTestResult result = HitTestResult();
      binding.hitTestInView(result, getCenter(finder), this.view.viewId);

      blocked = result.path.any((entry) {
        final target = entry.target;
        return (target is RenderIgnorePointer && target.ignoring) ||
            (target is RenderAbsorbPointer && target.absorbing);
      });

      frames++;
    }
    await pumpAndSettle();
  }
}

abstract class Resolutions {
  static const Map<String, Size> all = {
    'mobile portrait small (360x640)': Size(360, 640),
    'mobile portrait (390x844)': Size(390, 844),
    'mobile landscape (844x390)': Size(844, 390),
    'tablet portrait (768x1024)': Size(768, 1024),
    'tablet landscape (1024x768)': Size(1024, 768),
    'desktop (1280x720)': Size(1280, 720),
    'desktop large (1920x1080)': Size(1920, 1080),
  };

  static ValueVariant<MapEntry<String, Size>> variant() => ValueVariant(all.entries.toSet());
}
