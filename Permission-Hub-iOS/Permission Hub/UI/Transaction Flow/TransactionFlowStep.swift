//
//  TransactionFlowStep.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

enum TransactionFlowStep: CaseIterable {
    case overview, personalDetails, identityDocuments, paymentMethod, bankVerification, newsletter, legalTerms, terms1, terms2, finalOverView

    var title: String? {
        switch self {
        case .overview:
            return nil
        case .personalDetails:
            return "Personal details"
        case .identityDocuments:
            return "Identity documents"
        case .paymentMethod:
            return "Payment method"
        case .bankVerification:
            return "Bank verification"
        case .newsletter:
            return "Newsletter subscription"
        case .legalTerms:
            return "Legal terms"
        case .terms1:
            return ""
        case .terms2:
            return ""
        case .finalOverView:
            return ""
        }
    }

    func viewController(withTransaction transaction: TransactionNotification) -> UIViewController {
        switch self {
        case .overview:
            return TransactionOverviewViewController(transaction: transaction, isFinal: false)
        case .personalDetails:
            return TransactionPersonalDetailsViewController(transaction: transaction)
        case .identityDocuments:
            return TransactionIdentityDocumentViewController(transaction: transaction)
        case .paymentMethod:
            return TransactionPaymentMethodViewController(transaction: transaction)
        case .bankVerification:
            return TransactionPaymentConfirmationViewController(transaction: transaction)
        case .newsletter:
            return TransactionNewsletterSubscriptionViewController(transaction: transaction)
        case .legalTerms:
            return TransactionLegalTermsViewController(transaction: transaction)
        case .terms1:
            return TransactionTermsViewController(fileName: "terms1")
        case .terms2:
            return TransactionTermsViewController(fileName: "terms2")
        case .finalOverView:
            return TransactionOverviewViewController(transaction: transaction, isFinal: true)
        }
    }
}
