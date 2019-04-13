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
            return TransactionOverViewViewController(transaction: transaction)
        case .personalDetails:
            return PHTableViewController(title: title, items: [PHTableViewViewCellType]())
        case .identityDocuments:
            return PHTableViewController(title: title, items: [PHTableViewViewCellType]())
        case .paymentMethod:
            return PHTableViewController(title: title, items: [PHTableViewViewCellType]())
        case .bankVerification:
            return PHTableViewController(title: title, items: [PHTableViewViewCellType]())
        }
    }
}
