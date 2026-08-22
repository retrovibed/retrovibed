import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/library/api.dart' as api;
import 'package:retrovibed/media/player.control.resume.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

class _FakePlatformPlayer extends PlatformPlayer {
  _FakePlatformPlayer() : super(configuration: const PlayerConfiguration());
}

Player _fakePlayer() => Player(platformPlayer: _FakePlatformPlayer());

void main() {
  group('PlayerControlResume overflow', () {
    testWidgets('renders a long title without horizontal overflow', (tester) async {
      final player = _fakePlayer();
      addTearDown(player.dispose);
      final current = api.Known(
        id: 'a',
        description: 'A Very Long Episode Title That Would Not Fit On One Line Without Wrapping Or Truncating',
      );

      await tester.pumpApp(
        PlayerControlResume(player, current, ValueNotifier<bool>(true)),
        physicalSize: const Size(320, 200),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(find.byType(PlayerControlResume), findsOneWidget);
    });

    testWidgets('renders a short title without overflow', (tester) async {
      final player = _fakePlayer();
      addTearDown(player.dispose);
      final current = api.Known(id: 'a', description: 'Short Title');

      await tester.pumpApp(
        PlayerControlResume(player, current, ValueNotifier<bool>(true)),
        physicalSize: const Size(320, 200),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(find.textContaining('Short Title'), findsOneWidget);
    });
  });

  group('PlayerControlResume cursor', () {
    testWidgets('shows click cursor over the control', (tester) async {
      final player = _fakePlayer();
      addTearDown(player.dispose);
      final current = api.Known(id: 'a', description: 'Title');

      await tester.pumpApp(
        PlayerControlResume(player, current, ValueNotifier<bool>(true)),
      );
      await tester.pumpAndSettle();

      expect(
        tester.resolvedCursorAt(find.byType(PlayerControlResume)),
        SystemMouseCursors.click,
      );
    });
  });
}
