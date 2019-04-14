//
//  TransactionPaymentConfirmationViewController.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionPaymentConfirmationViewController: PHTableViewController {

    // MARK: - Private properties

    private var transaction: TransactionNotification

    // MARK: - Initialization

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        super.init(title: "", items: [
            .notification(
                type: .verification,
                text: "This company is verified"),
            .description(
                image: nil,
                date: transaction.date,
                title: "Directly from bank account",
                description: "The Payment Services Directive 2 requires your express consent in order to allow the Third Party Provider (Stripe) to access your bank account information currently stored by your bank. This will allow Stripe to make payment on your behalf. You can withdraw your consent anytime."),
            .notification(
                type: .verification,
                text: "This Third Party provider is verified"),
            .plugin(
                image: UIImage(named: "stripe"),
                text: "By clicking the green button below you allow Stripe to access the bank account information from your bank in order to make recurrent payments on your behalf."),
            .selection(options: ["Allow access to my bank information"])
            ])
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        tableView.allowsSelection = false
    }
}
