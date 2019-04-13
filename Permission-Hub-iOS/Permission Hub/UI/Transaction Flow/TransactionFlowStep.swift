//
//  TransactionFlowStep.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

enum TransactionFlowStep: CaseIterable {
    case overview, personalDetails, identityDocuments, paymentMethod, bankVerification

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
        }
    }

    func viewController(withTransaction transaction: TransactionNotification) -> UIViewController {
        switch self {
        case .overview:
            return TransactionOverviewViewController(transaction: transaction)
        case .personalDetails:
            return TransactionPersonalDetailsViewController(transaction: transaction)
        case .identityDocuments:
            return TransactionIdentityDocumentViewController(transaction: transaction)
        case .paymentMethod:
            return TransactionPaymentMethodViewController(transaction: transaction)
        case .bankVerification:
            return PHTableViewController(title: title, items: [PHTableViewViewCellType]())
        }
    }
}
