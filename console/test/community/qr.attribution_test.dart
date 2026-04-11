import 'package:flutter_test/flutter_test.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/community/qr.attribution.dart';
import 'package:retrovibed/community/community.pb.dart';

final _community = Community(
  domain: 'testdomain',
  description: 'A test community',
  createdAt: '2024-01-15T14:30:00Z',
);

final _resolutions = Resolutions.variant();

void main() {
  group('QRAttribution', () {
    group('resolutions', () {
      testWidgets('renders without overflow', (tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          QRAttribution(community: _community),
        );
        await tester.pumpAndSettle();

        expect(find.byType(QrImageView), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });
  });
}
