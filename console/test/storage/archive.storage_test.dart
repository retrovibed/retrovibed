import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/meta/meta.authn.pb.dart' as authn_meta;
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/storage/archive.storage.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/quotas.dart' as quotas;
import 'package:flutter/material.dart';

Future<quotas.QuotaFindResponse> mockQuota(
  String sku, {
  List<httpx.Option> options = const [],
}) async {
  return quotas.QuotaFindResponse(
    quota: quotas.Quota(consumed: fixnum.Int64(0)),
  );
}

Future<authn_meta.Authed> mockAuthSsh() async {
  return authn_meta.Authed(
    profiles: [authn_meta.Authn(token: "test-token")],
  );
}

void main() {
  group('ArchiveStorage', () {
    testWidgets('renders within 384 width', (WidgetTester tester) async {
      await tester.pumpApp(
        authn.Authenticated(
          Center(
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 384, maxHeight: 384),
              child: ArchiveStorage(quota: mockQuota),
            ),
          ),
          apissh: mockAuthSsh,
          apicurrent: (token) async => authn.Session(token: token),
        ),
      );
      await tester.pumpAndSettle();

      final RenderBox box = tester.renderObject(find.byType(ArchiveStorage));
      expect(box.size.width, 384.0);
      expect(box.size.height, 91.0);
      expect(tester.takeException(), isNull);
    });
  });
}
